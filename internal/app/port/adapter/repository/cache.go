package repository

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
)

type RateCache interface {
	GetRate(ctx context.Context, key redis.RateKey) (*redis.Rate, error)
	PutRate(ctx context.Context, key redis.RateKey, value redis.Rate) error
}
