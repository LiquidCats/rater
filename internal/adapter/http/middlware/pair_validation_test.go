package middlware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/internal/adapter/http/middlware"
	"github.com/LiquidCats/rater/internal/adapter/http/server"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestPairValidationMiddleware_Handle(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		pair entity.CurrencyPairString
	}{
		{
			name: "empty",
			pair: " ",
			msg:  "empty pair",
		},
		{
			name: "short",
			pair: "BTC_US",
			msg:  "should be 7 characters long at least",
		},
		{
			name: "lower",
			pair: "btc_usd",
			msg:  "should be uppercase",
		},
		{
			name: "not allowed",
			pair: "BTC_EUR",
			msg:  "not allowed",
		},
	}

	l := zerolog.New(zerolog.NewTestWriter(t))

	allowed := []entity.CurrencyPairString{
		"BTC_USD",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := server.NewRouter(&l)

			handler := middlware.NewPairValidation(allowed)

			router.GET("/rate/:pair", handler.Handle)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/rate/", tt.pair), nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

			assert.Contains(t, w.Body.String(), tt.msg)
		})
	}
}
