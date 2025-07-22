package configs

import (
	"fmt"

	"github.com/LiquidCats/rater/pkg/docker"
)

type DB struct {
	Driver   string `envconfig:"DRIVER" default:"postgres"`
	Host     string `envconfig:"HOST"`
	Port     string `envconfig:"PORT"`
	Database string `envconfig:"DATABASE"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
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
