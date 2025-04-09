package coingate

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
)

type Repository struct {
	cfg configs.CoinGateConfig
}

func NewRepository(cfg configs.CoinGateConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (c *Repository) GetRate(ctx context.Context, pair entity.Pair) (big.Float, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		c.cfg.URL,
		pair.From.ToLower(),
		pair.To.ToLower(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: error making http request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not read response body")
	}

	v, _, err := big.ParseFloat(string(resBody), 10, 0, big.ToNearestEven) // nolint:mnd
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not parse response")
	}

	return *v, nil
}

func (c *Repository) Name() entity.ProviderName {
	return "coingate"
}
