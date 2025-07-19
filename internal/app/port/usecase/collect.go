package usecase

import (
	"context"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type CollectRateUseCase interface {
	CollectRate(ctx context.Context, pair entity.Symbol) error
}
