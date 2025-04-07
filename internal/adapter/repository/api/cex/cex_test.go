package cex_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/cex"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/cex/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"

	"github.com/stretchr/testify/require"
)

var response = `{"ok":"ok","data":{"BTC-USD":{"bestBid":"66786.9","bestAsk":"66842.7","bestBidChange":"-3353.0","bestBidChangePercentage":"-4.78","bestAskChange":"-3310.2","bestAskChangePercentage":"-4.71","low":"65465.0","high":"71221.5","volume30d":"231.72466546","lastTradeDateISO":"2024-04-12T19:21:07.446Z","volume":"9.27624555","quoteVolume":"639276.97927628","lastTradeVolume":"0.00000458","volumeUSD":"639276.97","last":"66842.7","lastTradePrice":"66830.7","priceChange":"-3310.2","priceChangePercentage":"-4.71"}}}`

func TestRepository_Get(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		var req data.ApiRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		require.Contains(t, req.Pairs, "BTC-USD")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer ts.Close()

	repo := cex.NewRepository(configs.CexConfig{URL: ts.URL})

	ctx := context.Background()

	res, err := repo.GetRate(ctx, entity.Pair{
		From: "BTC",
		To:   "USD",
	})

	require.NoError(t, err)

	require.Equal(t, "66830.7", res.String())
}
