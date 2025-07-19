package repository

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/cache/redis"
)

type RateCache interface {
	GetRate(ctx context.Context, key string) (*redis.Rate, error)
	PutRate(ctx context.Context, key string, value redis.Rate) error
}
