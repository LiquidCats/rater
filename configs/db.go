package configs

import (
	"fmt"

	"github.com/LiquidCats/rater/pkg/docker"
)

type DB struct {
	Driver   string `yaml:"driver" envconfig:"DRIVER" default:"postgres"`
	Host     string `yaml:"host" envconfig:"HOST"`
	Port     string `yaml:"port" envconfig:"PORT"`
	Database string `yaml:"database" envconfig:"DATABASE"`
	User     string `yaml:"user" envconfig:"USER"`
	Password string `yaml:"password" envconfig:"PASSWORD"`
}

func (d *DB) ToDSN() string {
	pwd, _ := docker.GetSecret(d.Password)

	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s?sslmode=disable",
		d.Driver,
		d.User,
		pwd,
		d.Host,
		d.Port,
		d.Database,
	)
}
