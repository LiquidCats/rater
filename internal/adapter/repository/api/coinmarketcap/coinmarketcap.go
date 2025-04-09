package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinmarketcap/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
)

type Repository struct {
	cfg configs.CoinMarketCapConfig
}

func NewReposiotry(cfg configs.CoinMarketCapConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (big.Float, error) {
	url := fmt.Sprintf(
		"%s?amount=1&symbol=%s&convert=%s",
		r.cfg.URL,
		pair.From.ToLower(),
		pair.To.ToLower(),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not create request")
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", string(r.cfg.Secret)) //nolint:canonicalheader

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: error making http request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	var response data.APIResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not parse response")
	}

	quotePrice, ok := response.Data.Quote[pair.To.ToUpper().String()]
	if !ok {
		return big.Float{}, errors.Wrap(err, "repo: could not parse response")
	}

	return *big.NewFloat(quotePrice.Price), nil
}
