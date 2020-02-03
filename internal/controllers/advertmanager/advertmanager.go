package advertmanager

import (
	"github.com/jackc/pgx"
)
//TODO: заворачивать статус код
type AdvertManager struct {
	ConnPool   *pgx.ConnPool
	HTTPStatus uint
}
