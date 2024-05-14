package port

import (
	"context"
	"math/big"
	"rater/internal/app/domain/types"
)

type RateRepository interface {
	Get(ctx context.Context, quote types.QuoteCurrency, base types.BaseCurrency) (*big.Float, error)
}

type NamedRateRepository interface {
	Name() string
}
