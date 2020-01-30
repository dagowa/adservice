package store

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/joeshaw/envdecode"
)

type redisConfig struct {
	Network string `env:"NETWORK,default=tcp"`
	Address string `env:"ADDRESS,default=localhost:6379"`
}

// A RedisConn is implementation of Redis connection
// and its configuration
type RedisConn struct {
	config *redisConfig
	Conn   *redis.Pool
}

func newREDISConfig(path string) (*redisConfig, error) {
	config := new(redisConfig)
	if err := envdecode.StrictDecode(config); err != nil {
		return nil, errors.New("Cannot decode redis configuration; err: " + err.Error())
	}
	if path != "" {
		if err := config.setEnvs(path); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (rc *redisConfig) setEnvs(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return errors.New("Cannot open file with configuration; err: " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return errors.New("Cannot read file; err: " + err.Error())
	}
	if err := json.Unmarshal(byteValue, &rc); err != nil {
		return errors.New("Cannot unmarshal JSON; err: " + err.Error())
	}
	return nil
}

// NewREDISConnection initializes REDIS connection configuration with configPath
// Pass configPath == "" for default configuration
func NewREDISConnection(configPath string) (*RedisConn, error) {
	rc := new(RedisConn)
	config, err := newREDISConfig(configPath)
	if err != nil {
		return nil, errors.New("Cannot set configuration; err: " + err.Error())
	}
	rc.config = config

	rc.Conn = &redis.Pool{
		MaxIdle:   100,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial(rc.config.Network, rc.config.Address)
			if err != nil {
				return nil, errors.New("Cannot init redis connection; err: " + err.Error())
			}
			return conn, nil
		},
	}
	if rc.Conn == nil {
		return nil, errors.New("Connection failed")
	}
	return rc, nil
}
