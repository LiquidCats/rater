package entity

import (
	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	"github.com/shopspring/decimal"
)

type Rate struct {
	Pair     Pair
	Price    decimal.Decimal
	Provider ProviderName
}

func NewRate(rate postgres.Rate) Rate {
	return Rate{
		Pair:     Symbol(rate.Pair).ToPair(),
		Price:    rate.Price,
		Provider: ProviderName(rate.Provider),
	}
}
