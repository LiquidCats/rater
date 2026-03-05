package coingate

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	fullURL := fmt.Sprintf(
		"%s/%s/%s",
		c.cfg.URL,
		pair.From.ToLower(),
		pair.To.ToLower(),
	)

	parsedURL, err := url.ParseRequestURI(fullURL)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: incorrect request url")
	}

	if parsedURL.Scheme != "https" && parsedURL.Scheme != "http" {
		return decimal.Zero, eris.New("repo: unsupported URL scheme")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not create request")
	}

	res, err := http.DefaultClient.Do(req) //nolint:gosec
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
