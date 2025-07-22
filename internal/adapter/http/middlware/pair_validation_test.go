package middlware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	http2 "github.com/LiquidCats/rater/internal/adapter/http"
	"github.com/LiquidCats/rater/internal/adapter/http/middlware"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/LiquidCats/rater/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPairValidationMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		pair entity.Symbol
		mock func(t *testing.T, symbol entity.Symbol) repository.PairDatabase
	}{
		{
			name: "empty",
			pair: " ",
			msg:  "empty pair",
			mock: func(t *testing.T, _ entity.Symbol) repository.PairDatabase {
				pairDB := mocks.NewPairDatabase(t)
				pairDB.On("CountPairs", mock.Anything).Once().Return(int64(1), nil)

				return pairDB
			},
		},
		{
			name: "short",
			pair: "BTC_US",
			mock: func(t *testing.T, _ entity.Symbol) repository.PairDatabase {
				pairDB := mocks.NewPairDatabase(t)
				pairDB.On("CountPairs", mock.Anything).Once().Return(int64(1), nil)

				return pairDB
			},
			msg: "should be 7 characters long at least",
		},
		{
			name: "lower",
			pair: "btc_usd",
			msg:  "should be uppercase",
			mock: func(t *testing.T, _ entity.Symbol) repository.PairDatabase {
				pairDB := mocks.NewPairDatabase(t)
				pairDB.On("CountPairs", mock.Anything).Once().Return(int64(1), nil)

				return pairDB
			},
		},
		{
			name: "not allowed",
			pair: "BTC_EUR",
			msg:  "not allowed",
			mock: func(t *testing.T, symbol entity.Symbol) repository.PairDatabase {
				pairDB := mocks.NewPairDatabase(t)
				pairDB.On("CountPairs", mock.Anything).Once().Return(int64(1), nil)
				pairDB.On("HasPair", mock.Anything, symbol.ToUpper().String()).Once().Return(false, nil)

				return pairDB
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := http2.NewRouter()

			pairDB := tt.mock(t, tt.pair)

			handler := middlware.NewPairValidation(pairDB)

			router.GET("/rate/:pair", handler.Handle)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/rate/", tt.pair), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

			assert.Contains(t, w.Body.String(), tt.msg)
		})
	}
}
