package repository

import (
	"context"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/shopspring/decimal"
)

type RateAPI interface {
	GetRate(ctx context.Context, pair entity.Pair) (decimal.Decimal, error)
}
