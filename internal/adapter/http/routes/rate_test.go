package routes_test

import (
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/http/routes"
	"github.com/LiquidCats/rater/internal/adapter/http/server"
	"github.com/LiquidCats/rater/internal/adapter/repository/api"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_Handle(t *testing.T) {
	l := zerolog.New(zerolog.NewTestWriter(t))

	rateCache := mocks.NewRateCache(t)
	rateApi := mocks.NewRateApi(t)

	matcher := mock.MatchedBy(func(pair entity.Pair) bool {
		return pair.From == "BTC" && pair.To == "USD"
	})

	rateCache.On("GetRate", mock.Anything, matcher).Once().Return(nil, nil)
	rateCache.On("PutRate", mock.Anything, entity.Rate{
		Pair: entity.Pair{
			From: "BTC",
			To:   "USD",
		},
		Price:    *big.NewFloat(25000.77733333),
		Provider: "test",
	}, time.Minute*5).Once().Return(nil)

	rateApi.On("GetRate", mock.Anything, matcher).Once().Return(*big.NewFloat(25000.77733333), nil)

	useCase := usecase.NewRateUsecase(rateCache, api.Registry{
		"test": rateApi,
	})
	handler := routes.NewRateHandler(useCase)

	router := server.NewRouter(&l)

	router.GET("/rate/:pair", handler.Handle)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/rate/", "BTC_USD"), nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"data":{"pair":"BTC_USD","price":{}},"status":"success"}`, w.Body.String())
}
