package data

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type CoinGeckoID string

func (i CoinGeckoID) String() string {
	return string(i)
}

var mapping = map[entity.CurrencyISO]CoinGeckoID{ // nolint:gochecknoglobals
	"BTC": "bitcoin",
	"ETH": "ethereum",
}

func GetCoinID(iso entity.CurrencyISO) (CoinGeckoID, bool) {
	id, ok := mapping[iso.ToUpper()]

	return id, ok
}
