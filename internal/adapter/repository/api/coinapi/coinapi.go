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
	"github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/rotisserie/eris"
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
		return decimal.Zero, eris.Wrap(err, "repo: could not create request")
	}

	secret, err := r.cfg.GetSecret()
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not get secret")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CoinAPI-Key", string(secret)) //nolint:canonicalheader

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

	var resp data.APIResponse

	if err := decoder.Decode(&resp); err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not unmarshal response")
	}

	return decimal.NewFromFloat(resp.Rate), nil
}
