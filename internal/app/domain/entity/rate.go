package entity

import (
	"github.com/shopspring/decimal"
)

type Rate struct {
	Pair     Pair            `json:"pair"`
	Price    decimal.Decimal `json:"price"`
	Provider ProviderName    `json:"provider"`
}
