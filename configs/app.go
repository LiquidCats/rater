package configs

import (
	"github.com/rs/zerolog"
)

type AppConfig struct {
	LogLevel zerolog.Level `envconfig:"LOG_LEVEL" default:"info"`
}
