package postgresql

import (
	"time"

	"github.com/jackc/pgx"
)

// A Connection is implementation of PostgreSQL connection
// and its configuration
type Connection struct {
	config *Config
	pool   *pgx.ConnPool
}

// NewConnection initializes ppool connection configuration with configPath
func NewConnection(cfg *Config) (*Connection, error) {
	conn := Connection{
		config: cfg,
	}

	var runtimeParams map[string]string
	runtimeParams = make(map[string]string)
	runtimeParams["application_name"] = "adservice"

	connConfig := pgx.ConnConfig{
		User:              conn.config.User,
		Password:          conn.config.Password,
		Host:              conn.config.Host,
		Port:              conn.config.Port,
		Database:          conn.config.Dbname,
		TLSConfig:         nil,
		UseFallbackTLS:    conn.config.UseFallbackTLS,
		FallbackTLSConfig: nil,
		RuntimeParams:     runtimeParams,
	}
	pgxconfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: conn.config.MaxConnections,
		AcquireTimeout: time.Duration(conn.config.AcquireTimeout) * time.Second,
	}
	ppool, err := pgx.NewConnPool(pgxconfig)
	if err != nil {
		return nil, err
	}
	conn.pool = ppool
	return &conn, nil
}

// Close ends the use of a connection pool
func (c *Connection) Close() {
	c.pool.Close()
}

// Pool returns the postgres connection pool object
func (c *Connection) Pool() *pgx.ConnPool {
	return c.pool
}
