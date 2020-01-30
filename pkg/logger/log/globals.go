package log

import (
	"github.com/rs/zerolog/log"

	"github.com/dagowa/adservice/pkg/logger"
)

var (
	Logger logger.Logger = log.Logger

	Fatal = log.Fatal
)
