package controllers

import (
	"github.com/dagowa/adservice/internal/controllers/advertmanager"
	"github.com/jackc/pgx"
)

type Controllers struct {
	ConnPool *pgx.ConnPool
}

func NewControllerSet(pp *pgx.ConnPool) Controllers {
	return Controllers{
		ConnPool: pp,
	}
}

func (c *Controllers) AdvertManager() *advertmanager.AdvertManager {
	return &advertmanager.AdvertManager{ConnPool: c.ConnPool}
}
