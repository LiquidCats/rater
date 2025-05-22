package configs

import (
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/go-playground/sensitive"
)

type Config struct {
	App   AppConfig           `yaml:"app" envconfig:"APP"`
	Redis RedisConfig         `yaml:"redis" envconfig:"REDIS"`
	HTTP  graceful.HttpConfig `yaml:"http" envconfig:"HTTP"`

	CoinGate      CoinGateConfig      `yaml:"coingate" envconfig:"COIN_GATE"`
	Cex           CexConfig           `yaml:"cex" envconfig:"CEX"`
	CoinApi       CoinApiConfig       `yaml:"coin_api" envconfig:"COIN_API"` // nolint:revive
	CoinMarketCap CoinMarketCapConfig `yaml:"coin_market_cap" envconfig:"COIN_MARKET_CAP"`
	CoinGecko     CoinGeckoConfig     `yaml:"coin_gecko" envconfig:"COIN_GECKO"`
}

type AppConfig struct {
	LogLevel zerolog.Level `envconfig:"LOG_LEVEL" default:"info"`

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
	URL        string           `yaml:"url" envconfig:"URL"`
	Secret     sensitive.String `yaml:"secret" envconfig:"SECRET"`
	SecretFile string           `yaml:"secret_file" envconfig:"SECRET_FILE"`
}

func (c CoinApiConfig) GetSecret() (sensitive.String, error) {
	if len(c.Secret) > 0 {
		return c.Secret, nil
	}

	return getSecretFromFile(c.SecretFile)
}

type CoinMarketCapConfig struct {
	URL        string           `yaml:"url" envconfig:"URL"`
	Secret     sensitive.String `yaml:"secret" envconfig:"SECRET"`
	SecretFile string           `yaml:"secret_file" envconfig:"SECRET_FILE"`
}

func (c CoinMarketCapConfig) GetSecret() (sensitive.String, error) {
	if len(c.Secret) > 0 {
		return c.Secret, nil
	}

	return getSecretFromFile(c.SecretFile)
}

type CoinGeckoConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}

type RedisConfig struct {
	Address  string
	Password sensitive.String
	DB       int
}

func getSecretFromFile(filePath string) (sensitive.String, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "coinMarketCapConfig: couldn't read secret file")
	}

	return sensitive.String(b), nil
}
