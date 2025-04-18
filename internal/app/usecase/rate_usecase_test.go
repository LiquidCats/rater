package usecase_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExchange_Get(t *testing.T) {
	tests := []struct {
		name   string
		before func(t *testing.T) *usecase.RateUsecase
	}{
		{
			name: "rate from provider",
			before: func(t *testing.T) *usecase.RateUsecase {
				rateCache := mocks.NewRateCache(t)
				rateAPI := mocks.NewRateAPI(t)

				matcher := mock.MatchedBy(func(pair entity.Pair) bool {
					return pair.From == "USD" && pair.To == "BTC"
				})

				rateCache.On("GetRate", mock.Anything, matcher).Once().Return(nil, nil)
				rateCache.On("PutRate", mock.Anything, entity.Rate{
					Pair: entity.Pair{
						From: "USD",
						To:   "BTC",
					},
					Price:    decimal.NewFromFloat(25000.77733333),
					Provider: "test",
				}, time.Second*5).Once().Return(nil)

				rateAPI.On("GetRate", mock.Anything, matcher).Once().Return(decimal.NewFromFloat(25000.77733333), nil)

				providers := api.Registry{
					"test": rateAPI,
				}

				return usecase.NewRateUsecase(rateCache, providers)
			},
		},
		{
			name: "rate from cache",
			before: func(t *testing.T) *usecase.RateUsecase {
				rateCache := mocks.NewRateCache(t)
				rateAPI := mocks.NewRateAPI(t)

				matcher := mock.MatchedBy(func(pair entity.Pair) bool {
					return pair.From == "USD" && pair.To == "BTC"
				})

				rateCache.On("GetRate", mock.Anything, matcher).Once().Return(&entity.Rate{
					Pair: entity.Pair{
						From: "USD",
						To:   "BTC",
					},
					Price:    decimal.NewFromFloat(25000.77733333),
					Provider: "test",
				}, nil)

				providers := api.Registry{
					"test": rateAPI,
				}

				return usecase.NewRateUsecase(rateCache, providers)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			ctx := t.Context()

			uc := tt.before(t)

			// act
			rate, err := uc.GetRate(ctx, entity.Pair{
				From: "USD",
				To:   "BTC",
			})

			// assert
			require.NoError(t, err)
			require.NotNil(t, rate)

			assert.Equal(t, "BTC", string(rate.Pair.To))
			assert.Equal(t, "USD", string(rate.Pair.From))
			assert.Equal(t, "25000.77733333", rate.Price.String())
		})
	}
}
