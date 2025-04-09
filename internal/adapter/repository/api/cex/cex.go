package cex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/cex/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/pkg/errors"
)

type Repository struct {
	cfg configs.CexConfig
}

func NewRepository(cfg configs.CexConfig) *Repository {
	return &Repository{cfg: cfg}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (big.Float, error) {
	body := data.APIRequest{Pairs: []string{pair.Join("-")}}

	bodyByres, err := json.Marshal(body)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: incorrect request body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.cfg.URL, bytes.NewBuffer(bodyByres))
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

	var resp data.APIResponse

	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not unmarshal response")
	}

	tickerInfo, ok := resp.Data[body.Pairs[0]]
	if !ok {
		return big.Float{}, errors.Wrap(err, "repo: could not unmarshal response")
	}

	v, _, err := big.ParseFloat(tickerInfo.LastTradePrice, 10, 0, big.ToNearestEven) // nolint:mnd
	if err != nil {
		return big.Float{}, errors.Wrap(err, "repo: could not parse value")
	}

	return *v, nil
}
