package configs

import (
	"os"
	"strconv"
	"strings"
)

func Load() Config {
	baseCurrencies := strings.Split(os.Getenv("RATER_BASE_CURRENCIES"), ",")
	quoteCurrencies := strings.Split(os.Getenv("RATER_QUOTE_CURRENCIES"), ",")

	redisDB, err := strconv.Atoi(os.Getenv("RATER_REDIS_DB"))
	if nil != err {
		redisDB = 0
	}

	return Config{
		Port:            os.Getenv("RATER_PORT"),
		QuoteCurrencies: quoteCurrencies,
		BaseCurrencies:  baseCurrencies,
		Redis: RedisConfig{
			Address:  os.Getenv("RATER_REDIS_ADDRESS"),
			Password: os.Getenv("RATER_REDIS_PASSWORD"),
			DB:       redisDB,
		},
	}
}
