package repository

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type RateCache interface {
	GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error)
	PutRate(ctx context.Context, rate entity.Rate, expire time.Duration) error
}
