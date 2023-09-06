package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"rater/internal/adapter/api/dto"
	"testing"
	"time"
)

var quote = "USD"
var base = "BTC"
var host = fmt.Sprintf("http://dev:8080/v1/rate/%s/%s", base, quote)
var key = fmt.Sprintf("rater:rate:base:%s:quote:%s", base, quote)

type response struct {
	Data struct {
		Price string `json:"price"`
		Base  string `json:"base"`
		Quote string `json:"quote"`
	} `json:"data"`
	Status string `json:"status"`
}

func getRate(ctx context.Context) (*response, error) {
	req, err := http.NewRequest(http.MethodGet, host, nil)
	if nil != err {
		return nil, err
	}

	req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var rateResponse *response

	err = json.Unmarshal(body, &rateResponse)
	if nil != err {
		return nil, err
	}

	return rateResponse, nil
}

func TestGetRate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	redisClient := redis.NewClient(&redis.Options{
		Addr: "cache:6379",
		DB:   0,
	})

	redisClient.FlushAll(ctx)

	t.Run("first run", func(t *testing.T) {

		rateResponse, err := getRate(ctx)
		require.NoError(t, err)

		if err := redisClient.Exists(ctx, key).Err(); !errors.Is(err, redis.Nil) {
			require.NoError(t, err)
		}

		require.Equal(t, dto.StatusSuccess, rateResponse.Status)
		require.Equal(t, base, rateResponse.Data.Base)
		require.Equal(t, quote, rateResponse.Data.Quote)
		require.Equal(t, "29295.929694597355", rateResponse.Data.Price)
	})

	t.Run("cached", func(t *testing.T) {
		if err := redisClient.Exists(ctx, key).Err(); !errors.Is(err, redis.Nil) {
			require.NoError(t, err)
		}

		rateResponse, err := getRate(ctx)
		require.NoError(t, err)

		require.Equal(t, dto.StatusSuccess, rateResponse.Status)
		require.Equal(t, base, rateResponse.Data.Base)
		require.Equal(t, quote, rateResponse.Data.Quote)
		require.Equal(t, "29295.92969", rateResponse.Data.Price)
	})
}
