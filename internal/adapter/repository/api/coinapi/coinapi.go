package coinapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coinapi/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
)

type Repository struct {
	cfg configs.CoinApiConfig
}

func NewRepository(cfg configs.CoinApiConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (a *Repository) GetRate(ctx context.Context, pair entity.Pair) (big.Float, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		a.cfg.URL,
		pair.From.ToUpper(),
		pair.To.ToUpper(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CoinAPI-Key", string(a.cfg.Secret))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: error making http request")
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return big.Float{}, errors.Errorf("repo: error making http request: %s", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not read response body")
	}

	var resp data.ApiResponse

	if err = json.Unmarshal(resBody, &resp); nil != err {
		return big.Float{}, errors.Wrap(err, "repo: could not parse response")
	}

	return *big.NewFloat(resp.Rate), nil
}
