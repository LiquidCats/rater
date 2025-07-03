package cex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/LiquidCats/rater/configs"
	"github.com/LiquidCats/rater/internal/adapter/repository/api/cex/data"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/domain/errors"
	"github.com/rotisserie/eris"
	"github.com/shopspring/decimal"
)

type Repository struct {
	cfg configs.CexConfig
}

func NewRepository(cfg configs.CexConfig) *Repository {
	return &Repository{cfg: cfg}
}

func (r *Repository) GetRate(ctx context.Context, pair entity.Pair) (decimal.Decimal, error) {
	body := data.APIRequest{Pairs: []string{pair.Join("-")}}

	bodyByres, err := json.Marshal(body)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: incorrect request body")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.cfg.URL, bytes.NewBuffer(bodyByres))
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
		if err = decoder.Decode(&resBody); err != nil && eris.Is(err, io.EOF) {
			return decimal.Zero, eris.Wrap(err, "repo: could not decode response body")
		}

		return decimal.Zero, &errors.ProviderRequestFailedError{
			StatusCode: res.StatusCode,
			Body:       resBody,
		}
	}

	var resp data.APIResponse

	if err = decoder.Decode(&resp); err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not unmarshal response")
	}

	tickerInfo, ok := resp.Data[body.Pairs[0]]
	if !ok {
		return decimal.Zero, eris.Wrap(err, "repo: could not unmarshal response")
	}

	value, err := decimal.NewFromString(tickerInfo.LastTradePrice)
	if err != nil {
		return decimal.Zero, eris.Wrap(err, "repo: could not parse value")
	}

	return value, nil
}
