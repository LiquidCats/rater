package configs

type Config struct {
	Port            string `default:"8080"`
	BaseCurrencies  []string
	QuoteCurrencies []string
	Redis           RedisConfig
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}
