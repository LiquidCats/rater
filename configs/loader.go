package configs

import (
	"github.com/kelseyhightower/envconfig"
)

func Load(prefix string) (Config, error) {
	var cfg Config

	if err := envconfig.Process(prefix, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
