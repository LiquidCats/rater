package entity

import "math/big"

type Rate struct {
	Pair     Pair
	Price    big.Float
	Provider ProviderName
}
