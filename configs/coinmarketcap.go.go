package configs

import (
	"os"

	"github.com/LiquidCats/graceful"
	"github.com/go-playground/sensitive"
	"github.com/rotisserie/eris"
)

type Config struct {
	App     AppConfig           `yaml:"app" envconfig:"APP"`
	Redis   RedisConfig         `yaml:"redis" envconfig:"REDIS"`
	HTTP    graceful.HttpConfig `yaml:"http" envconfig:"HTTP"`
	Metrics graceful.HttpConfig `yaml:"metrics" envconfig:"METRICS"`
	DB      DB                  `yaml:"db" envconfig:"DB"`

	CoinGate      CoinGateConfig      `yaml:"coingate" envconfig:"COIN_GATE"`
	Cex           CexConfig           `yaml:"cex" envconfig:"CEX"`
	CoinApi       CoinApiConfig       `yaml:"coin_api" envconfig:"COIN_API"` // nolint:revive
	CoinMarketCap CoinMarketCapConfig `yaml:"coin_market_cap" envconfig:"COIN_MARKET_CAP"`
	CoinGecko     CoinGeckoConfig     `yaml:"coin_gecko" envconfig:"COIN_GECKO"`
}

func getSecretFromFile(filePath string) (sensitive.String, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return "", eris.Wrapf(err, "coinMarketCapConfig: couldn't read secret file")
	}

	return sensitive.String(b), nil
}
