package cron_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	"github.com/LiquidCats/rater/internal/adapter/scheduler/cron"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

func TestCollectRateTask_Run(t *testing.T) {
	logger := zerolog.New(zerolog.NewTestWriter(t))
	cfg := configs.AppConfig{
		CollectSchedule: "",
	}

	pairDB := mocks.NewPairDatabase(t)
	cache := mocks.NewCacheService(t)
	svc := mocks.NewRateService(t)

	pairs := []postgres.Pair{
		{
			Symbol:     "BTC_USD",
			BaseAsset:  "BTC",
			QuoteAsset: "USD",
			CreatedAt: pgtype.Timestamp{
				Time: time.Now(),
			},
		},
	}

	pair := entity.Pair{
		From:   "BTC",
		To:     "USD",
		Symbol: "BTC_USD",
	}

	rate := entity.Rate{
		Pair:     pair,
		Price:    decimal.RequireFromString("120435.78933333"),
		Provider: "testprovider",
	}
	pairDB.On("GetAllPairs", mock.Anything).Once().Return(pairs, nil)
	svc.On("Current", mock.Anything, rate.Pair).Once().Return(&rate, rate.Provider, nil)

	uc := usecase.NewRateUseCase(usecase.RateUseCaseDeps{
		Cache:   cache,
		Service: svc,
	})

	task := cron.NewCollectRateTask(
		&logger,
		cfg,
		pairDB,
		uc,
	)

	task.Run()
}
