package repository

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/database"
)

type RateDatabase interface {
	SaveRate(ctx context.Context, arg database.SaveRateParams) (database.Rate, error)
}
