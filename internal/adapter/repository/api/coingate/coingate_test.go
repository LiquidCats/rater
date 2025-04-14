package coingate_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coingate"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const response = `29295.929694597355`

func TestCoinGateRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := coingate.NewRepository(configs.CoinGateConfig{
		URL: ts.URL,
	})

	rate, err := repo.GetRate(t.Context(), entity.Pair{
		From: "BTC",
		To:   "USD",
	})
	require.NoError(t, err)
	require.NotNil(t, rate)

	require.Equal(t, "29295.929694597355", rate.String())
}
