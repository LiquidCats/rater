package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"rater/internal/app/domain/types"
)

type Repository struct {
	url    string
	secret string
}

type price struct {
	Price float64
}

type conversion struct {
	ID     string
	Amount uint
	Name   string
	Symbol string
	Quote  map[string]price
}

type apiResponse struct {
	Data conversion
}

func NewReposiotry() *Repository {
	url := os.Getenv("RATER_COINMARKETCAP_URL")
	secret := os.Getenv("RATER_COINMARKETCAP_SECRET")

	return &Repository{
		url:    url,
		secret: secret,
	}
}

func (r *Repository) Get(ctx context.Context, quote types.QuoteCurrency, base types.BaseCurrency) (*big.Float, error) {
	url := fmt.Sprintf(
		"%s?amount=1&symbol=%s&convert=%s",
		r.url,
		base.Upper(),
		quote.Lower(),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("repo: could not create request: %s\n", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", r.secret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo: error making api request: %s\n", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(res.Body)

	var response apiResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("repo: could not parse response: %w", err)
	}

	quotePrice, ok := response.Data.Quote[quote.Upper()]
	if !ok {
		return nil, fmt.Errorf("repo: could not parse response: %w", err)
	}

	return big.NewFloat(quotePrice.Price), nil
}

func (r *Repository) Name() string {
	return "coinmarketcap"
}
