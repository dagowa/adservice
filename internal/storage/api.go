package storage

import (
	"github.com/dagowa/adservice/internal/storage/postgresql"
	"github.com/dagowa/adservice/internal/storage/redis"
)

type storage struct {
}

// New defines new storage interface
func New() *storage {
	return &storage{}
}

func (storage) NewPostgreSQLConn(cfg *postgresql.Config) (*postgresql.Connection, error) {
	conn, err := postgresql.NewConnection(cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (storage) NewRedisConn(cfg *redis.Config) (*redis.Connection, error) {
	conn, err := redis.NewConnection(cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
