package service

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type CacheService interface {
	GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error)
	PutRate(ctx context.Context, rate entity.Rate) error
}

type RateService interface {
	Historical(ctx context.Context, pair entity.Pair, ts time.Time) (*entity.Rate, error)
	Current(ctx context.Context, pair entity.Pair) (*entity.Rate, entity.ProviderName, error)
}
