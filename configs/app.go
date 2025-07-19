package configs

import (
	"github.com/rs/zerolog"
)

const AppName = "rater"

type AppConfig struct {
	LogLevel zerolog.Level `envconfig:"LOG_LEVEL" default:"info"`

	CollectSchedule string `envconfig:"COLLECT_SCHEDULE" default:"*/5 * * * *"`
}
