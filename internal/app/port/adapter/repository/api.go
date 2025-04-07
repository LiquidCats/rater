package repository

import (
	"context"
	"math/big"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type RateApi interface {
	GetRate(ctx context.Context, pair entity.Pair) (big.Float, error)
}
