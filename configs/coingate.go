package configs

type CoinGateConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}
