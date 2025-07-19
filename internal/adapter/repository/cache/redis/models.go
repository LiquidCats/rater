package redis

import "github.com/shopspring/decimal"

type Rate struct {
	From     string          `json:"from"`
	To       string          `json:"to"`
	Price    decimal.Decimal `json:"price"`
	Provider string          `json:"provider"`
}
