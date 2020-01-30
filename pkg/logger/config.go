package logger

type Config struct {
	Level     string `env:"LOGGER_LEVEL"`
	Timestamp bool   `env:"LOGGER_TIMESTAMP"`
	Caller    bool   `env:"LOGGER_CALLER"`
	Pretty    bool   `env:"LOGGER_PRETTY"`
}
