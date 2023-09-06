package api

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type CoinGateRepository struct {
	url string
}

func NewCoinGateRepository(url string) *CoinGateRepository {
	return &CoinGateRepository{url: url}
}

func (c *CoinGateRepository) Get(ctx context.Context, quote, base string) (*big.Float, error) {
	url := fmt.Sprintf("%s/%s/%s", c.url, strings.ToUpper(base), strings.ToUpper(quote))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("repo: could not create request: %s\n", err)
	}

	req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo: error making api request: %s\n", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("repo: could not read response body: %s\n", err)
	}

	v, _, err := big.ParseFloat(string(resBody), 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("repo: could not parse response: %s\n", err)
	}

	return v, nil
}

func (c *CoinGateRepository) Name() string {
	return "coindate"
}
