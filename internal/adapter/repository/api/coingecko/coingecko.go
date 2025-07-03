package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coingecko/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type Repository struct {
	cfg configs.CoinGeckoConfig
}

func NewRepository(cfg configs.CoinGeckoConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (decimal.Decimal, error) {
	// https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=8
	id, ok := data.GetCoinID(pair.From)
	if !ok {
		return decimal.Zero, eris.New("repo: cant find coingecko id")
	}

	url := fmt.Sprintf(
		"%s?ids=%s&vs_currencies=%s&precision=8",
		r.cfg.URL,
		id,
		pair.To.ToLower(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: error making http request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode >= http.StatusBadRequest {
		var resBody string
		if err := decoder.Decode(&resBody); err != nil && err != io.EOF {
			return decimal.Zero, eris.Wrap(err, "repo: could not decode response body")
		}

		return decimal.Zero, &errors.ErrProviderRequestFailed{
			StatusCode: res.StatusCode,
			Body:       resBody,
		}
	}

	var response data.APIResponse

	if err = decoder.Decode(&response); err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not parse response")
	}

	rateBase, ok := response[id.String()]
	if !ok {
		return decimal.Zero, eris.New("repo: could not get base rate from response")
	}

	rateQuote, ok := rateBase.(map[string]interface{})
	if !ok {
		return decimal.Zero, eris.New("repo: could not get quoted rate from response")
	}

	val, ok := rateQuote[pair.To.ToLower().String()]
	if !ok {
		return decimal.Zero, eris.New("repo: could not get rate value from response")
	}

	floatVal, ok := val.(float64)
	if !ok {
		return decimal.Zero, eris.New("repo: could not get float value from response")
	}

	return decimal.NewFromFloat(floatVal), nil
}
