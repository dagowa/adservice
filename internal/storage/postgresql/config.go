package postgresql

// Config of postgres connection
type Config struct {
	Driver         string `env:"DRIVER_NAME,default=postgres"`
	User           string `env:"USER,default=postgres"`
	Password       string `env:"PASSWORD,default=v9bnVv31n"`
	Host           string `env:"HOST,default=localhost"`
	Port           uint16 `env:"PORT,default=5432"`
	Dbname         string `env:"DBNAME,default=adservice"`
	UseFallbackTLS bool   `env:"USE_FALLBACK_TLS,default=false"`
	MaxConnections int    `env:"MAX_CONNECTIONS,default=8"`
	AcquireTimeout int    `env:"ACQUIRE_TIMEOUT,default=2"`
}
