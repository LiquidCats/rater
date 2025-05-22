package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	http2 "github.com/LiquidCats/rater/internal/adapter/http"
	"github.com/LiquidCats/rater/internal/adapter/http/routes"
	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_Handle(t *testing.T) {
	rateCache := mocks.NewRateCache(t)
	rateAPI := mocks.NewRateAPI(t)

	matcher := mock.MatchedBy(func(pair entity.Pair) bool {
		return pair.From == "BTC" && pair.To == "USD"
	})

	rateCache.On("GetRate", mock.Anything, matcher).Once().Return(nil, nil)
	rateCache.On("PutRate", mock.Anything, entity.Rate{
		Pair: entity.Pair{
			From: "BTC",
			To:   "USD",
		},
		Price:    decimal.NewFromFloat(25000.77733333),
		Provider: "test",
	}, time.Second*5).Once().Return(nil)

	rateAPI.On("GetRate", mock.Anything, matcher).Once().Return(decimal.NewFromFloat(25000.77733333), nil)

	useCase := usecase.NewRateUsecase(rateCache, api.Registry{
		"test": rateAPI,
	})
	handler := routes.NewRateHandler(useCase)

	router := http2.NewRouter()

	router.GET("/rate/:pair", handler.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/rate/", "BTC_USD"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"data":{"pair":"BTC_USD","price":"25000.77733333"},"status":"success"}`, w.Body.String())
}
