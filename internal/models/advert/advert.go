package advert

import (
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type Advert struct {
	ID          int       `json:"id,omitempty"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	Date        time.Time `json:"date,omitempty"`
	Description *string   `json:"description,omitepmty"`
	Gallery     *[]string `json:"gallery"`
}

func New(title string, price int, date time.Time, description string, gallery *[]string) *Advert {
	return &Advert{
		Title:       title,
		Price:       price,
		Date:        date,
		Description: &description,
		Gallery:     gallery,
	}
}

func GetSingle(p *pgx.Conn, IDadvert int, isDateReq, isDescrReq, isFullGalleryReq bool) (*Advert, error) {
	selectStm := "SELECT a.id_advert, a.title, a.price, "
	joinStm := ""
	if isDateReq {
		selectStm += "a.description, "
	}
	if isDescrReq {
		selectStm += "a.date, "
	}
	if isFullGalleryReq {
		selectStm += "ARRAY_AGG(g.index), ARRAY_AGG(g.photo), "
		joinStm += "JOIN photo_gallery pg ON pg.id_advert=a.id_advert"
	}
	selectStm = strings.TrimSuffix(selectStm, ", ")

	var adv Advert

	row := p.QueryRow(selectStm+
		" FROM advert a"+
		joinStm+
		"WHERE a.id_advert = $1 "+
		"GROUP BY a.id_advert", IDadvert)

	if isDateReq {
		if isDescrReq {
			if isFullGalleryReq {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description, &adv.Gallery); err != nil {
					return nil, err
				}
			} else {
				if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
					&adv.Date, &adv.Description); err != nil {
					return nil, err
				}
			}
		} else {
			if err := row.Scan(&adv.ID, &adv.Title, &adv.Price,
				&adv.Date); err != nil {
				return nil, err
			}
		}
	} else {
		if err := row.Scan(&adv.ID, &adv.Title, &adv.Price); err != nil {
			return nil, err
		}
	}
	return &adv, nil
}
