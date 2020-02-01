package advert

import (
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type Advert struct {
	ID          int      `json:"id,omitempty"`
	Title       string   `json:"title"`
	Price       int      `json:"price"`
	Date        string   `json:"date,omitempty"`
	Description *string  `json:"description,omitepmty"`
	Gallery     *[]Photo `json:"gallery"`
}

type Photo struct {
	Index int    `json:"index"`
	Link  string `json:"photo"`
}

func New(title string, price int, date time.Time, description string, gallery *[]Photo) *Advert {
	return &Advert{
		Title:       title,
		Price:       price,
		Date:        date.Format("2006-01-02 15:04:05"),
		Description: &description,
		Gallery:     gallery,
	}
}

func GetSingle(p *pgx.Conn, IDadvert int, isDateReq, isDescrReq, isFullGalleryReq bool) (*Advert, error) {
	selectStm := "SELECT a.id_advert, a.title, a.price, "
	joinStm := "JOIN photo_gallery pg ON pg.id_advert=a.id_advert"
	whereStm := "WHERE a.id_advert = $1 "
	if isDateReq {
		selectStm += "a.description, "
	}
	if isDescrReq {
		selectStm += "a.date, "
	}
	if isFullGalleryReq {
		selectStm += "ARRAY_AGG(pg.index), ARRAY_AGG(pg.photo), "
	} else {
		selectStm += "gp.index, gp.photo, "
		whereStm += "AND pg.index = 0 "
	}
	selectStm = strings.TrimSuffix(selectStm, ", ")

	var adv Advert
	var indexList []int
	var photoList []string

	var mainPhotoIndex int
	var mainPhotoLink string

	row := p.QueryRow(selectStm+
		" FROM advert a"+
		joinStm+
		whereStm+
		"GROUP BY a.id_advert", IDadvert)

	if isDateReq {

		if isDescrReq {

			if isFullGalleryReq {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description, &indexList, &photoList); err != nil {
					return nil, err
				}
				var gallery []Photo
				for i := range indexList {
					gallery = append(gallery, Photo{
						Index: indexList[i],
						Link:  photoList[i]})
				}
				adv.Gallery = &gallery

			} else {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description, &mainPhotoIndex, &mainPhotoLink); err != nil {
					return nil, err
				}
				var gallery []Photo
				gallery = append(gallery, Photo{
					Index: mainPhotoIndex,
					Link:  mainPhotoLink})
				adv.Gallery = &gallery
			}

		} else {
			if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
				&adv.Date, &mainPhotoIndex, &mainPhotoLink); err != nil {
				return nil, err
			}
			var gallery []Photo
			gallery = append(gallery, Photo{
				Index: mainPhotoIndex,
				Link:  mainPhotoLink})
			adv.Gallery = &gallery
		}

	} else {
		if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
			&mainPhotoIndex, &mainPhotoLink); err != nil {
			return nil, err
		}
		var gallery []Photo
		gallery = append(gallery, Photo{
			Index: mainPhotoIndex,
			Link:  mainPhotoLink})
		adv.Gallery = &gallery
	}
	return &adv, nil
}

func (a *Advert) Add(p *pgx.ConnPool) (int, error) {
	addAdvQuery := "INSERT INTO public.advert( " +
		"title, price, description) " +
		"VALUES ($1, $2, $3) RETURNING id_advert"
	addGalleryQuery := "INSERT INTO public.photo_gallery( " +
		"id_advert, index, photo) " +
		"VALUES ($1, $2, $3);"
	tx, err := p.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	if err := tx.QueryRow(addAdvQuery,
		a.Title, a.Price, a.Description).Scan(&id); err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}

	var rowsAdded int64

	gallery := *(a.Gallery)
	for j := range gallery {
		ctag, err := tx.Exec(addGalleryQuery,
			id, gallery[j].Index, gallery[j].Link)

		if err != nil {
			if err := tx.Rollback(); err != nil {
				return 0, err
			}
			return 0, err
		}
		rowsAdded += ctag.RowsAffected()

	}
	if rowsAdded != int64(len(gallery)) {
		if err := tx.Rollback(); err != nil {
			return 0, err
		}
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func Delete(p *pgx.ConnPool, id int) error {
	tx, err := p.Begin()
	if err != nil {
		return err
	}
	deleteAdvQuery := "DELETE FROM public.photo_gallery " +
		"WHERE id_advert = $1; " +
		"DELETE FROM public.advert " +
		"WHERE id_advert = $2;"
	_, err = tx.Exec(deleteAdvQuery, id, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func AddBatch(p *pgx.ConnPool, advList *[]Advert) error {
	addAdvQuery := "INSERT INTO public.advert( " +
		"title, price, description) " +
		"VALUES ($1, $2, $3) RETURNING id_advert"
	addGalleryQuery := "INSERT INTO public.photo_gallery( " +
		"id_advert, index, photo) " +
		"VALUES ($1, $2, $3);"

	tx, err := p.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Prepare("add_adv_batch", addAdvQuery)
	if err != nil {
		return err
	}
	_, err = tx.Prepare("add_gallery_batch", addGalleryQuery)
	if err != nil {
		return err
	}

	list := *advList
	for i := range list {
		var id int
		if err := tx.QueryRow("add_adv_batch",
			list[i].Title, list[i].Price, list[i].Description).Scan(&id); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}

		var rowsAdded int64

		gallery := *(list[i].Gallery)
		for j := range gallery {
			ctag, err := tx.Exec("add_gallery_batch",
				id, gallery[j].Index, gallery[j].Link)

			if err != nil {
				if err := tx.Rollback(); err != nil {
					return err
				}
				return err
			}
			rowsAdded += ctag.RowsAffected()
		}

		if rowsAdded != int64(len(gallery)) {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
