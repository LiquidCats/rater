package usecase_test

import (
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRateUseCase_GetRate(t *testing.T) {
	tests := []struct {
		name   string
		ts     time.Time
		before func(t *testing.T) *usecase.RateUseCase
	}{
		{
			name: "get rate now",
			ts:   time.Now(),
			before: func(t *testing.T) *usecase.RateUseCase {
				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}

				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				service.On("Current", mock.Anything, rate.Pair).Once().Return(&rate, rate.Provider, nil)

				cache.On("GetRate", mock.Anything, rate.Pair).Once().Return(nil, nil)
				cache.On("PutRate", mock.Anything, rate).Once().Return(nil)

				return usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})
			},
		},
		{
			name: "get rate historical",
			ts:   time.Now(),
			before: func(t *testing.T) *usecase.RateUseCase {
				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				cache.On("GetRate", mock.Anything, rate.Pair).Once().Return(&rate, nil)

				return usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})
			},
		},
		{
			name: "get historical",
			ts:   time.Date(2021, 4, 15, 11, 42, 31, 0, time.UTC),
			before: func(t *testing.T) *usecase.RateUseCase {
				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				ts := time.Date(2021, 4, 15, 11, 40, 0, 0, time.UTC)
				service.On("Historical", mock.Anything, rate.Pair, ts).Once().Return(&rate, nil)

				return usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			ctx := t.Context()

			uc := tt.before(t)

			// act
			rate, err := uc.GetRate(ctx, "BTC_USD", tt.ts)

			// assert
			require.NoError(t, err)
			require.NotNil(t, rate)

			assert.Equal(t, "BTC", string(rate.Pair.From))
			assert.Equal(t, "USD", string(rate.Pair.To))
			assert.Equal(t, "25000.77733333", rate.Price.String())
		})
	}
}

func TestRateUseCase_CollectRate(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name   string
		before func(t *testing.T) *usecase.RateUseCase
		after  func(t *testing.T, err error)
	}{
		{
			name: "happy path",
			before: func(t *testing.T) *usecase.RateUseCase {
				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}

				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				service.On("Current", mock.Anything, rate.Pair).Once().Return(&rate, rate.Provider, nil)

				return usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})
			},
			after: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "negative",
			before: func(t *testing.T) *usecase.RateUseCase {
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				service.On("Current", mock.Anything, entity.Symbol("BTC_USD").ToPair()).
					Once().
					Return(
						nil,
						entity.ProviderName("test"),
						eris.New("test error"),
					)

				return usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})
			},
			after: func(t *testing.T, err error) {
				require.Error(t, err)
				assert.True(t, eris.Is(err, eris.New("test error")))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := test.before(t)
			require.NotNil(t, uc)

			err := uc.CollectRate(ctx, "BTC_USD")

			test.after(t, err)
		})
	}
}
