package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"rater/internal/adapter/api/dto"
	"rater/internal/adapter/api/middlware"
	"rater/internal/adapter/api/routes"
	"rater/internal/adapter/logger"
	"rater/internal/adapter/repository/api/coinapi"
	"rater/internal/adapter/repository/cache/memory"
	"rater/internal/app/domain/types"
	"rater/internal/app/usecase"
	"rater/internal/port"
	"testing"
	"time"
)

var quote types.CurrencyNameString = "USD"
var base types.CurrencyNameString = "BTC"
var uri = fmt.Sprintf("/v1/rate/%s/%s", base, quote)
var key = fmt.Sprintf("rate:base:%s:quote:%s", base, quote)

type response struct {
	Data struct {
		Price string `json:"price"`
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"data"`
	Status string `json:"status"`
}

func makeApp() (*gin.Engine, port.CacheRepository) {
	log := logger.NewNilLogger()

	cache := memory.NewCacheRepository()

	rateUsecase := usecase.NewRateUsecase(log, cache)

	rateUsecase.SetAdapter(coinapi.NewRepository())

	rateHandler := routes.NewRateHandler(rateUsecase)

	baseCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("base", []types.CurrencyNameString{base})
	quoteCurrencyMiddleware := middlware.NewCurrencyParamMiddleware("quote", []types.CurrencyNameString{quote})

	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	v1Router := router.Group("/v1")
	v1Router.GET(
		"/rate/:base/:quote",
		baseCurrencyMiddleware.Handle,
		quoteCurrencyMiddleware.Handle,
		rateHandler.GetRate,
	)

	return router, cache
}

func requestRate(ctx context.Context, url string) (*response, error) {
	req, err := http.NewRequest(http.MethodGet, url+uri, nil)
	if nil != err {
		return nil, err
	}

	req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)

	var rateResponse *response

	err = json.Unmarshal(body, &rateResponse)
	if nil != err {
		return nil, err
	}

	return rateResponse, nil
}

func TestGetRate(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"time":"2023-07-24T11:31:56.0000000Z","asset_id_base":"BTC","asset_id_quote":"USD","rate":29295.929694597355}`))
	}))
	defer api.Close()

	err := os.Setenv("RATER_COINAPI_URL", api.URL)
	require.NoError(t, err)

	router, cache := makeApp()

	server := httptest.NewServer(router)
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	t.Run("first run", func(t *testing.T) {
		if cache.Has(ctx, key) {
			require.Failf(t, "cache should not exist", "cache key: %s", key)
		}

		rateResponse, err := requestRate(ctx, server.URL)
		require.NoError(t, err)

		if !cache.Has(ctx, key) {
			require.Failf(t, "cache should exist", "cache key: %s", key)
		}

		assert.Equal(t, dto.StatusSuccess, rateResponse.Status)
		assert.Equal(t, string(base), rateResponse.Data.Base)
		assert.Equal(t, string(quote), rateResponse.Data.Quote)
		assert.Equal(t, "29295.929694597355", rateResponse.Data.Price)
	})

	t.Run("cached", func(t *testing.T) {
		if !cache.Has(ctx, key) {
			require.Failf(t, "cache should exist", "cache key: %s", key)
		}

		rateResponse, err := requestRate(ctx, server.URL)
		require.NoError(t, err)

		require.Equal(t, dto.StatusSuccess, rateResponse.Status)
		require.Equal(t, string(base), rateResponse.Data.Base)
		require.Equal(t, string(quote), rateResponse.Data.Quote)
		require.Equal(t, "29295.92969", rateResponse.Data.Price)
	})
}
