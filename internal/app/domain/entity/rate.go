package entity

import "math/big"

type Rate struct {
	Quote string
	Base  string
	Price *big.Float
}
