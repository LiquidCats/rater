package repository

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
)

type RateDatabase interface {
	GetRate(ctx context.Context, arg postgres.GetRateParams) (postgres.Rate, error)
	HasRate(ctx context.Context, arg postgres.HasRateParams) (bool, error)
	SaveRate(ctx context.Context, arg postgres.SaveRateParams) (postgres.Rate, error)
}

type PairDatabase interface {
	GetAllPairs(ctx context.Context) ([]postgres.Pair, error)
	GetPair(ctx context.Context, symbol string) (postgres.Pair, error)
	CountPairs(ctx context.Context) (int64, error)
	HasPair(ctx context.Context, symbol string) (bool, error)
}

type ProviderDatabase interface {
	GetAllProviders(ctx context.Context) ([]postgres.Provider, error)
}
