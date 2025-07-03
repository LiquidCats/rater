package configs

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/rs/zerolog"
)

type AppConfig struct {
	LogLevel zerolog.Level `envconfig:"LOG_LEVEL" default:"info"`

	Pairs []entity.CurrencyPairString
}
