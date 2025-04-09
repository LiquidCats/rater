package configs

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"

	"github.com/go-playground/sensitive"
)

type Config struct {
	App   AppConfig   `yaml:"app" envconfig:"APP"`
	Redis RedisConfig `yaml:"redis" envconfig:"REDIS"`

	CoinGate      CoinGateConfig      `yaml:"coingate" envconfig:"COINGATE"`
	Cex           CexConfig           `yaml:"cex" envconfig:"CEX"`
	CoinApi       CoinApiConfig       `yaml:"coin_api" envconfig:"COIN_API"` // nolint:revive
	CoinMarketCap CoinMarketCapConfig `yaml:"coin_market_cap" envconfig:"COIN_MARKET_CAP"`
	CoinGecko     CoinGeckoConfig     `yaml:"coin_gecko" envconfig:"COIN_GECKO"`
}

type AppConfig struct {
	Port  string `default:"8080"`
	Pairs []entity.CurrencyPairString
}

type CoinGateConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type CexConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type CoinApiConfig struct { // nolint:revive
	URL    string           `yaml:"url" envconfig:"URL"`
	Secret sensitive.String `yaml:"secret" envconfig:"SECRET"`
}

type CoinMarketCapConfig struct {
	URL    string           `yaml:"url" envconfig:"URL"`
	Secret sensitive.String `yaml:"secret" envconfig:"SECRET"`
}

type CoinGeckoConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type RedisConfig struct {
	Address  string
	Password sensitive.String
	DB       int
}
