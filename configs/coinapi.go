package configs

import "github.com/go-playground/sensitive"

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
