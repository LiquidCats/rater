package cex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
)

type Repository struct {
	url string
}

type requestBody struct {
	Pairs []string `json:"pairs"`
}

type responseBody struct {
	Ok   string                        `json:"ok"`
	Data map[string]responseBodyTicker `json:"data"`
}

type responseBodyTicker struct {
	LastTradePrice string `json:"lastTradePrice"`
}

func NewRepository() *Repository {
	url := os.Getenv("RATER_CEX_URL") // https://trade.cex.io/api/spot/rest-public/get_ticker

	return &Repository{url: url}
}

func (a *Repository) Get(ctx context.Context, quote, base string) (*big.Float, error) {
	pair := fmt.Sprintf("%s-%s", strings.ToUpper(base), strings.ToUpper(quote))

	body := requestBody{Pairs: []string{pair}}

	bodyByres, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("repo: incorrect request body: %s\n", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.url, bytes.NewBuffer(bodyByres))
	if err != nil {
		return nil, fmt.Errorf("repo: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo: error making api request: %s\n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("repo: could not read response body: %s\n", err)
	}

	var resp responseBody

	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("repo: could not unmarshal response: %s\n", err)
	}

	tickerInfo, ok := resp.Data[pair]
	if !ok {
		return nil, fmt.Errorf("repo: could not unmarshal response: %s\n", err)
	}

	v, _, err := big.ParseFloat(tickerInfo.LastTradePrice, 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("repo: could not parse value: %s\n", err)
	}

	return v, nil
}

func (a *Repository) Name() string {
	return "cex"
}
