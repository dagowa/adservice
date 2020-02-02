package controllers

import (
	"github.com/dagowa/adservice/internal/controllers/advertmanager"
	"github.com/jackc/pgx"
)

// NewAdvertManager returns advert manager interface
func NewAdvertManager(psql *pgx.ConnPool) *advertmanager.AdvertManager {
	return &advertmanager.AdvertManager{
		ConnPool: psql,
	}
}
