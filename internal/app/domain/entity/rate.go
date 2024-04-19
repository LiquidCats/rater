package entity

import "math/big"

type Rate struct {
	Quote    string     `json:"quote"`
	Base     string     `json:"base"`
	Price    *big.Float `json:"price"`
	Provider string     `json:"provider"`
}
