package redis

// Config of redis connection
type Config struct {
	Network string `env:"NETWORK,default=tcp"`
	Address string `env:"ADDRESS,default=localhost:6379"`
}
