package advertmanager

import (
	"strings"

	"github.com/dagowa/adservice/internal/models/advert"
	"github.com/jackc/pgx"
)

// AdvertManager implements basic manipulations with adverts
type AdvertManager struct {
	ConnPool     *pgx.ConnPool
	sortCriteria *SortCriteria
}

// GetOne gets one advert with basic fileds and extra fileds, described in
// RequirementFields object
func (am *AdvertManager) GetOne(id_advert int, rf *RequirementFileds) (*advert.Advert, error) {
	selectStm := "SELECT a.id_advert, a.title, a.price, "
	joinStm := "JOIN photo_gallery pg ON pg.id_advert=a.id_advert"
	whereStm := "WHERE a.id_advert = $1 "
	if rf.Date {
		selectStm += "a.description, "
	}
	if rf.Description {
		selectStm += "a.date, "
	}
	if rf.Gallery {
		selectStm += "ARRAY_AGG(pg.index), ARRAY_AGG(pg.photo), "
	} else {
		selectStm += "gp.index, gp.photo, "
		whereStm += "AND pg.index = 0 "
	}
	selectStm = strings.TrimSuffix(selectStm, ", ")

	var adv advert.Advert
	var indexList []int
	var photoList []string

	var mainPhotoIndex int
	var mainPhotoLink string

	row := am.ConnPool.QueryRow(selectStm+
		" FROM advert a"+
		joinStm+
		whereStm+
		"GROUP BY a.id_advert", id_advert)

	if rf.Date {

		if rf.Description {

			if rf.Gallery {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description, &indexList, &photoList); err != nil {
					return nil, err
				}
				var gallery []advert.Photo
				for i := range indexList {
					gallery = append(gallery, advert.Photo{
						Index: indexList[i],
						Link:  photoList[i]})
				}
				adv.Gallery = &gallery

			} else {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description, &mainPhotoIndex, &mainPhotoLink); err != nil {
					return nil, err
				}
				var gallery []advert.Photo
				gallery = append(gallery, advert.Photo{
					Index: mainPhotoIndex,
					Link:  mainPhotoLink})
				adv.Gallery = &gallery
			}

		} else {
			if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
				&adv.Date, &mainPhotoIndex, &mainPhotoLink); err != nil {
				return nil, err
			}
			var gallery []advert.Photo
			gallery = append(gallery, advert.Photo{
				Index: mainPhotoIndex,
				Link:  mainPhotoLink})
			adv.Gallery = &gallery
		}

	} else {
		if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
			&mainPhotoIndex, &mainPhotoLink); err != nil {
			return nil, err
		}
		var gallery []advert.Photo
		gallery = append(gallery, advert.Photo{
			Index: mainPhotoIndex,
			Link:  mainPhotoLink})
		adv.Gallery = &gallery
	}
	return &adv, nil
}

// GetBatch gets a batch of adverts with pre-defined filter criteria
func (am *AdvertManager) GetBatch(sc *SortCriteria, pageNum int, pageSize int) (*[]advert.Advert, error) {
	limitOffsetRow := "LIMIT $1 OFFSET $2"
	orderByRow := makeOrderByRow(sc)
	query := "SELECT a.id_advert, a.title, a.price, pg.photo " +
		"FROM advert a " +
		"JOIN photo_gallery pg ON a.id_advert=pg.id_advert " +
		orderByRow + limitOffsetRow

	tx, err := am.ConnPool.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(query, pageNum, pageSize)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	defer rows.Close()

	var advList []advert.Advert
	for rows.Next() {
		var adv advert.Advert
		var photo advert.Photo
		if err := rows.Scan(&adv.ID, &adv.Title, &adv.Price, &photo.Link); err != nil {
			if err := tx.Rollback(); err != nil {
				return nil, err
			}
			return nil, err
		}
		*(adv.Gallery) = append(*(adv.Gallery), photo)
		advList = append(advList, adv)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &advList, nil
}

// AddOne adds one advert to database
func (am *AdvertManager) AddOne(a *advert.Advert) (int, error) {
	addAdvQuery := "INSERT INTO public.advert( " +
		"title, price, description) " +
		"VALUES ($1, $2, $3) RETURNING id_advert"
	addGalleryQuery := "INSERT INTO public.photo_gallery( " +
		"id_advert, index, photo) " +
		"VALUES ($1, $2, $3);"
	tx, err := am.ConnPool.Begin()
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

// AddBatch adds a  batch of adverts to database
func (am *AdvertManager) AddBatch(aList *[]advert.Advert) error {
	addAdvQuery := "INSERT INTO public.advert( " +
		"title, price, description) " +
		"VALUES ($1, $2, $3) RETURNING id_advert"
	addGalleryQuery := "INSERT INTO public.photo_gallery( " +
		"id_advert, index, photo) " +
		"VALUES ($1, $2, $3);"

	tx, err := am.ConnPool.Begin()
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

	list := *aList
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

// Delete designed to delete advert with certain id
func (am *AdvertManager) Delete(id int) error {
	tx, err := am.ConnPool.Begin()
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

func makeOrderByRow(sc *SortCriteria) string {
	orderByRow := "ORDER BY "

	isSortComplex := false

	if sc.priceAsc {
		orderByRow += "a.price "
		isSortComplex = true
	} else {
		isSortComplex = true
		orderByRow += "a.price DESC "
	}

	if sc.dateAsv {
		if isSortComplex {
			orderByRow += ", a.date "
		}
		orderByRow += "a.date "
	} else {
		if isSortComplex {
			orderByRow += ", a.date DESC "
		}
		orderByRow += "a.date DESC "
	}
	return orderByRow
}
