package server

// Config is ...
type Config struct {
	Host        string `env:"SERVICE_HOST, required"`
	Port        int    `env:"SERVICE_PORT, required"`
	LogRequests bool   `env:"SERVICE_LOG_REQUESTS,default=true"`
}
