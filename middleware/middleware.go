package middleware

import (
	"github.com/jackc/pgx"
)

type Middleware struct {
	ConnPool *pgx.ConnPool
}
