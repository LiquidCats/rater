package coingecko

import "rater/internal/app/domain/types"

var mapping = map[types.BaseCurrency]string{
	"BTC": "bitcoin",
	"ETH": "ethereum",
}

func getCoinID(base types.BaseCurrency) string {
	return mapping[base]
}
