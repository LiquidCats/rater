package entity

import "math/big"

type Rate struct {
	Pair     Pair         `json:"pair"`
	Price    big.Float    `json:"price"`
	Provider ProviderName `json:"provider"`
}
