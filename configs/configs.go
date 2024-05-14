package configs

import "rater/internal/app/domain/types"

type Config struct {
	Port            string `default:"8080"`
	BaseCurrencies  []types.BaseCurrency
	QuoteCurrencies []types.QuoteCurrency
	Redis           RedisConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}
