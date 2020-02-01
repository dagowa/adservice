package advert

import "github.com/jackc/pgx"

type batchSearch struct {
}

func GetAdvertBatch() *batchSearch {
	return &batchSearch{}
}

func (b *batchSearch) PriceAscending(p *pgx.ConnPool, page_numb, size int) (*[]Advert, error) {
	query := "SELECT a.id_advert, a.title, a.price, pg.photo " +
		"FROM advert a " +
		"JOIN photo_gallery pg ON a.id_advert=pg.id_advert " +
		"ORDER BY a.price " +
		"LIMIT $1 OFFSET $2"
	adverts, err := b.getBatch(p, query, page_numb, size)
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (b *batchSearch) PriceDescending(p *pgx.ConnPool, page_numb, size int) (*[]Advert, error) {
	query := "SELECT a.id_advert, a.title, a.price, pg.photo " +
		"FROM advert a " +
		"JOIN photo_gallery pg ON a.id_advert=pg.id_advert " +
		"ORDER BY a.price DESC " +
		"LIMIT $1 OFFSET $2"
	adverts, err := b.getBatch(p, query, page_numb, size)
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (b *batchSearch) DateAscending(p *pgx.ConnPool, page_numb, size int) (*[]Advert, error) {
	query := "SELECT a.id_advert, a.title, a.price, pg.photo " +
		"FROM advert a " +
		"JOIN photo_gallery pg ON a.id_advert=pg.id_advert " +
		"ORDER BY a.date " +
		"LIMIT $1 OFFSET $2"
	adverts, err := b.getBatch(p, query, page_numb, size)
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (b *batchSearch) DateDescending(p *pgx.ConnPool, page_numb, size int) (*[]Advert, error) {
	query := "SELECT a.id_advert, a.title, a.price, pg.photo " +
		"FROM advert a " +
		"JOIN photo_gallery pg ON a.id_advert=pg.id_advert " +
		"ORDER BY a.date DESC " +
		"LIMIT $1 OFFSET $2"
	adverts, err := b.getBatch(p, query, page_numb, size)
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (batchSearch) getBatch(p *pgx.ConnPool, query string, page, size int) (*[]Advert, error) {
	tx, err := p.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query(query, page, size)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return nil, err
	}
	defer rows.Close()

	var advList []Advert
	for rows.Next() {
		var adv Advert
		var photo Photo
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
