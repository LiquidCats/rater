package configs

import (
	"github.com/LiquidCats/rater/pkg/docker"
	"github.com/go-playground/sensitive"
	"github.com/rotisserie/eris"
)

type CoinMarketCapConfig struct {
	URL        string           `yaml:"url" envconfig:"URL"`
	Secret     sensitive.String `yaml:"secret" envconfig:"SECRET"`
	SecretFile string           `yaml:"secret_file" envconfig:"SECRET_FILE"`
}

func (c CoinMarketCapConfig) GetSecret() (sensitive.String, error) {
	if len(c.Secret) > 0 {
		return c.Secret, nil
	}

	val, err := docker.GetSecret(c.SecretFile)
	if err != nil {
		return "", eris.Wrap(err, "coinMarketCapConfig: couldn't read secret file")
	}

	return sensitive.String(val), err
}
