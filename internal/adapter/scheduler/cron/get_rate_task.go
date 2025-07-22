package cron

import (
	"context"
	"time"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"
)

type CollectRateTask struct {
	logger *zerolog.Logger

	cfg     configs.AppConfig
	pairDB  repository.PairDatabase
	useCase *usecase.RateUseCase
}

func NewCollectRateTask(
	logger *zerolog.Logger,
	cfg configs.AppConfig,
	pairDB repository.PairDatabase,
	useCase *usecase.RateUseCase,
) *CollectRateTask {
	return &CollectRateTask{
		logger: logger,

		cfg:     cfg,
		pairDB:  pairDB,
		useCase: useCase,
	}
}

func (t *CollectRateTask) Spec() string {
	return t.cfg.CollectSchedule
}

func (t *CollectRateTask) Run() {
	t.logger.Info().Msg("start collect rate task")
	defer t.logger.Info().Msg("end collect rate task")

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute) //nolint:mnd
	defer cancel()

	pairs, err := t.pairDB.GetAllPairs(ctx)
	if err != nil {
		t.logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Msg("invalid pairs")
		return
	}

	for _, pair := range pairs {
		if err = t.useCase.CollectRate(ctx, entity.Symbol(pair.Symbol)); err != nil {
			t.logger.Error().
				Any("err", eris.ToJSON(err, true)).
				Str("symbol", pair.Symbol).
				Msg("invalid rate")
		}
	}
}
