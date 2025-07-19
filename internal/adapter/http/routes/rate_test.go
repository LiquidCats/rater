package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	http2 "github.com/LiquidCats/rater/internal/adapter/http"
	"github.com/LiquidCats/rater/internal/adapter/http/routes"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/metrics"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_Handle(t *testing.T) {
	tests := []struct {
		name   string
		ts     time.Time
		before func(t *testing.T) (*usecase.RateUseCase, metrics.ResponseTimeMetric)
	}{
		{
			name: "happy path",
			before: func(t *testing.T) (*usecase.RateUseCase, metrics.ResponseTimeMetric) {
				m := mocks.NewResponseTimeMetric(t)
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}

				service.On("Current", mock.Anything, rate.Pair).Once().Return(&rate, rate.Provider, nil)

				cache.On("GetRate", mock.Anything, rate.Pair).Once().Return(nil, nil)
				cache.On("PutRate", mock.Anything, rate).Once().Return(nil)

				m.On("Observe", "/rate/BTC_USD", mock.Anything).Once().Return()

				uc := usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})

				return uc, m
			},
		},
		{
			name: "from cache",
			before: func(t *testing.T) (*usecase.RateUseCase, metrics.ResponseTimeMetric) {
				m := mocks.NewResponseTimeMetric(t)
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}

				cache.On("GetRate", mock.Anything, rate.Pair).Once().Return(&rate, nil)
				m.On("Observe", "/rate/BTC_USD", mock.Anything).Once().Return()

				uc := usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})

				return uc, m
			},
		},
		{
			name: "get historical",
			ts:   time.Date(2021, 4, 15, 11, 42, 31, 0, time.UTC),
			before: func(t *testing.T) (*usecase.RateUseCase, metrics.ResponseTimeMetric) {
				m := mocks.NewResponseTimeMetric(t)
				cache := mocks.NewCacheService(t)
				service := mocks.NewRateService(t)

				rate := entity.Rate{
					Pair: entity.Pair{
						From:   "BTC",
						To:     "USD",
						Symbol: "BTC_USD",
					},
					Price:    decimal.RequireFromString("25000.77733333"),
					Provider: "test",
				}

				dt := time.Date(2021, 4, 15, 11, 42, 31, 0, time.UTC).
					Format(entity.DefaultFormat)
				ts := time.Date(2021, 4, 15, 11, 40, 0, 0, time.UTC)
				service.On("Historical", mock.Anything, rate.Pair, ts).Once().Return(&rate, nil)
				m.On("Observe", "/rate/BTC_USD/"+dt, mock.Anything).Once().Return()

				uc := usecase.NewRateUseCase(usecase.RateUseCaseDeps{
					Cache:   cache,
					Service: service,
				})

				return uc, m
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc, m := test.before(t)

			handler := routes.NewRateHandler(uc, routes.Metrics{
				ResponseTime: m,
			})

			router := http2.NewRouter()

			w := httptest.NewRecorder()
			var req *http.Request
			if test.ts.IsZero() {
				router.GET("/rate/:pair", handler.Handle)
				req, _ = http.NewRequest(http.MethodGet, fmt.Sprint("/rate/", "BTC_USD"), nil)
			} else {
				router.GET("/rate/:pair/*date", handler.Handle)
				req, _ = http.NewRequest(http.MethodGet, fmt.Sprint("/rate/BTC_USD/", test.ts.Format(entity.DefaultFormat)), nil)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.JSONEq(t, `{"data":{"pair":"BTC_USD","price":"25000.77733333"},"status":"success"}`, w.Body.String())
		})
	}
}
