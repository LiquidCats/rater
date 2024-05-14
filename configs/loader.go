package configs

import (
	"os"
	"rater/internal/app/domain/types"
	"strconv"
	"strings"
)

func Load() Config {

	cfg := Config{
		Port: os.Getenv("RATER_PORT"),
		Redis: RedisConfig{
			Address:  os.Getenv("RATER_REDIS_ADDRESS"),
			Password: os.Getenv("RATER_REDIS_PASSWORD"),
		},
	}

	loadRedisDB(&cfg)
	loadBaseCurrencies(&cfg)
	loadQuoteCurrencies(&cfg)

	return cfg
}

func loadRedisDB(cfg *Config) {
	redisDB, err := strconv.Atoi(os.Getenv("RATER_REDIS_DB"))
	if nil != err {
		redisDB = 0
	}

	cfg.Redis.DB = redisDB
}

func loadBaseCurrencies(cfg *Config) {
	baseCurrencies := strings.Split(os.Getenv("RATER_BASE_CURRENCIES"), ",")
	for _, currency := range baseCurrencies {
		cfg.BaseCurrencies = append(cfg.BaseCurrencies, types.BaseCurrency(currency))
	}
}

func loadQuoteCurrencies(cfg *Config) {
	quoteCurrencies := strings.Split(os.Getenv("RATER_QUOTE_CURRENCIES"), ",")
	for _, currency := range quoteCurrencies {
		cfg.QuoteCurrencies = append(cfg.QuoteCurrencies, types.QuoteCurrency(currency))
	}
}
