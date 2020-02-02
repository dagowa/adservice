package redis

import (
	"errors"

	"github.com/gomodule/redigo/redis"
)

// A Connection is implementation of Redis connection
// and its configuration
type Connection struct {
	config *Config
	pool   *redis.Pool
}

// NewConnection initializes REDIS connection configuration with configPath
func NewConnection(cfg *Config) (*Connection, error) {
	conn := Connection{
		config: cfg,
	}
	conn.pool = &redis.Pool{
		MaxIdle:   100,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(conn.config.Network, conn.config.Address)
			if err != nil {
				return nil, errors.New("Cannot init redis connection; err: " + err.Error())
			}
			return conn, nil
		},
	}
	if conn.pool == nil {
		return nil, errors.New("Connection failed")
	}
	return &conn, nil
}

// Close ends the use of a connection pool
func (c *Connection) Close() {
	c.pool.Close()
}

// Pool returns the redis connection pool object
func (c *Connection) Pool() *redis.Pool {
	return c.pool
}
