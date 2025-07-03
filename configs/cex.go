package configs

type CexConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}
