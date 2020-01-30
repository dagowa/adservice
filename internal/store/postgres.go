package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/jackc/pgx"

	"github.com/joeshaw/envdecode"
)

type psqlConfig struct {
	Driver         string `env:"DRIVER_NAME,default=postgres"`
	User           string `env:"USER,default=postgres"`
	Password       string `env:"PASSWORD,default=v9bnVv31n"`
	Host           string `env:"HOST,default=localhost"`
	Port           uint16 `env:"PORT,default=5432"`
	Dbname         string `env:"DBNAME,default=adservice"`
	MaxConnections int    `env:"MAX_CONNECTIONS,default=8"`
	AcquireTimeout int    `env:"ACQUIRE_TIMEOUT,default=2"`
}

// A PSQLConn is implementation of PostgreSQL connection
// and its configuration
type PSQLConn struct {
	config *psqlConfig
	Pool   *pgx.ConnPool
}

func newPSQLConfig(path string) (*psqlConfig, error) {
	config := new(psqlConfig)
	if err := envdecode.StrictDecode(config); err != nil {
		return nil, errors.New("Cannot decode postgres configuration; err: " + err.Error())
	}
	if path != "" {
		if err := config.setEnvs(path); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (pc *psqlConfig) setEnvs(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return errors.New("Cannot open file with configuration; err: " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.New("Cannot read file; err: " + err.Error())
	}
	if err := json.Unmarshal(byteValue, &pc); err != nil {
		return errors.New("Cannot unmarshal JSON; err: " + err.Error())
	}
	return nil
}

// NewPSQLConnection initializes ppool connection configuration with configPath
// Pass configPath == "" for default configuration
func NewPSQLConnection(configPath string) (*PSQLConn, error) {
	pc := new(PSQLConn)
	config, err := newPSQLConfig(configPath)
	if err != nil {
		return nil, errors.New("Cannot set configuration; err: " + err.Error())
	}
	pc.config = config

	var runtimeParams map[string]string
	runtimeParams = make(map[string]string)
	runtimeParams["application_name"] = "adservice"

	connConfig := pgx.ConnConfig{
		User:              pc.config.User,
		Password:          pc.config.Password,
		Host:              pc.config.Host,
		Port:              pc.config.Port,
		Database:          pc.config.Dbname,
		TLSConfig:         nil,
		UseFallbackTLS:    false,
		FallbackTLSConfig: nil,
		RuntimeParams:     runtimeParams,
	}
	pgxconfig := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		MaxConnections: pc.config.MaxConnections,
		AcquireTimeout: time.Duration(pc.config.AcquireTimeout) * time.Second,
	}
	ppool, err := pgx.NewConnPool(pgxconfig)
	if err != nil {
		return nil, err
	}
	pc.Pool = ppool
	return pc, nil
}
