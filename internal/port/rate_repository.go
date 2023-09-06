package port

import (
	"context"
	"math/big"
)

type RateRepository interface {
	Get(ctx context.Context, quote, base string) (*big.Float, error)
}

type NamedRateRepository interface {
	Name() string
}
