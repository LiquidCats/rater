package entity

type ProviderName string

const (
	ProviderNameCex           ProviderName = "cex"
	ProviderNameCoinApi       ProviderName = "coinapi" // nolint:revive
	ProviderNameCoinGate      ProviderName = "coingate"
	ProviderNameCoinGecko     ProviderName = "coingecko"
	ProviderNameCoinMarketCap ProviderName = "coinmarketcap"
)
