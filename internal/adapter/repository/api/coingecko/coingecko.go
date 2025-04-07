package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/coingecko/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
)

type Repository struct {
	cfg configs.CoinGeckoConfig
}

func NewRepository(cfg configs.CoinGeckoConfig) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (big.Float, error) {
	// https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=8
	id, ok := data.GetCoinID(pair.From)
	if !ok {
		return big.Float{}, errors.New("repo: cant find coingecko id")
	}

	url := fmt.Sprintf(
		"%s?ids=%s&vs_currencies=%s&precision=8",
		r.cfg.URL,
		id,
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	var response data.ApiResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not parse response")
	}

	rateBase, ok := response[id.String()]
	if !ok {
		return big.Float{}, errors.New("repo: could not get base rate from response")
	}

	rateQuote, ok := rateBase.(map[string]interface{})
	if !ok {
		return big.Float{}, errors.New("repo: could not get quoted rate from response")
	}

	val, ok := rateQuote[pair.To.ToLower().String()]
	if !ok {
		return big.Float{}, errors.New("repo: could not get rate value from reponse")
	}

	return *big.NewFloat(val.(float64)), nil
}
