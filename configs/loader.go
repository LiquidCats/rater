package configs

import (
	"github.com/kelseyhightower/envconfig"
)

func Load() (Config, error) {
	var cfg Config

	if err := envconfig.Process(AppName, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
