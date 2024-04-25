package entity

import (
	"math/big"
	"rater/internal/app/domain/types"
)

type Rate struct {
	Quote    types.QuoteCurrency `json:"quote"`
	Base     types.BaseCurrency  `json:"base"`
	Price    *big.Float          `json:"price"`
	Provider string              `json:"provider"`
}
