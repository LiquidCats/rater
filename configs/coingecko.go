package configs

type CoinGeckoConfig struct {
	URL string `yaml:"url" envconfig:"URL"`
}
