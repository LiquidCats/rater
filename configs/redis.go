package configs

import "github.com/go-playground/sensitive"

type RedisConfig struct {
	Address  string           `envconfig:"ADDRESS"`
	Password sensitive.String `envconfig:"PASSWORD"`
	DB       int              `envconfig:"DB" default:"0"`
	Protocol int              `envconfig:"PROTOCOL" default:"3"`
}
