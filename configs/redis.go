package configs

import "github.com/go-playground/sensitive"

type RedisConfig struct {
	Address  string
	Password sensitive.String
	DB       int
}
