package configs

type Config struct {
	Port            string
	BaseCurrencies  []string
	QuoteCurrencies []string
	CoinGateUrl     string
	CoinApiUrl      string
	CoinApiSecret   string
	Redis           RedisConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}
