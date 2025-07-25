package coingate

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type Repository struct {
	cfg configs.CoinGateConfig
}

func NewRepository(cfg configs.CoinGateConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (c *Repository) GetRate(ctx context.Context, pair entity.Pair) (decimal.Decimal, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		c.cfg.URL,
		pair.From.ToLower(),
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

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not read response body")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return decimal.Zero, &errors.ProviderRequestFailedError{
			StatusCode: res.StatusCode,
			Body:       string(data),
		}
	}

	value, err := decimal.NewFromString(string(data))
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not parse response")
	}

	return value, nil
}
