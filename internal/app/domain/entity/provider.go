package entity

type ProviderName string

func (p ProviderName) String() string {
	return string(p)
}

const (
	ProviderNameCex           ProviderName = "cex"
	ProviderNameCoinApi       ProviderName = "coinapi" // nolint:revive
	ProviderNameCoinGate      ProviderName = "coingate"
	ProviderNameCoinGecko     ProviderName = "coingecko"
	ProviderNameCoinMarketCap ProviderName = "coinmarketcap"
)

type Provider struct {
	Name ProviderName
}
