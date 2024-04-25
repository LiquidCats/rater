package coingate

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"rater/internal/app/domain/types"
)

type Repository struct {
	url string
}

func NewRepository() *Repository {
	url := os.Getenv("RATER_COINGATE_URL")

	return &Repository{url: url}
}

func (c *Repository) Get(ctx context.Context, quote types.QuoteCurrency, base types.BaseCurrency) (*big.Float, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		c.url,
		base.Lower(),
		quote.Lower(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	v, _, err := big.ParseFloat(string(resBody), 10, 0, big.ToNearestEven)
	if err != nil {
		return nil, fmt.Errorf("repo: could not parse response: %s\n", err)
	}

	return v, nil
}

func (c *Repository) Name() string {
	return "coingate"
}
