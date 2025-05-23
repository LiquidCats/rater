package coinapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinapi/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Repository struct {
	cfg configs.CoinApiConfig
}

func NewRepository(cfg configs.CoinApiConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (decimal.Decimal, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		r.cfg.URL,
		pair.From.ToUpper(),
		pair.To.ToUpper(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "repo: could not create request")
	}

	secret, err := r.cfg.GetSecret()
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "repo: could not get secret")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CoinAPI-Key", string(secret)) //nolint:canonicalheader

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "repo: error making http request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode >= http.StatusBadRequest {
		return decimal.Zero, errors.Errorf("repo: error making http request: %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return decimal.Zero, errors.Wrap(err, "repo: could not read response body")
	}

	var resp data.APIResponse

	if err = json.Unmarshal(resBody, &resp); nil != err {
		return decimal.Zero, errors.Wrap(err, "repo: could not parse response")
	}

	return decimal.NewFromFloat(resp.Rate), nil
}
