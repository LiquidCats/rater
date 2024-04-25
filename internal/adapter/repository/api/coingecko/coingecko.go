package coingecko

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
	url string
}

type apiResponse map[string]interface{}

func NewRepository() *Repository {
	url := os.Getenv("RATER_COINGECKO_API")

	return &Repository{
		url: url,
	}
}

func (r *Repository) Get(ctx context.Context, quote types.QuoteCurrency, base types.BaseCurrency) (*big.Float, error) {
	// https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd&precision=8
	url := fmt.Sprintf(
		"%s?ids=%s&vs_currencies=%s&precision=8",
		r.url,
		base.Transform(getCoinID),
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

	var response apiResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("repo: could not parse response: %w", err)
	}

	rateBase, ok := response[base.Transform(getCoinID)]
	if !ok {
		return nil, fmt.Errorf("repo: could not get base rate from response")
	}

	rateQuote, ok := rateBase.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("repo: could not get quoted rate from response")
	}

	val, ok := rateQuote[quote.Lower()]
	if !ok {
		return nil, fmt.Errorf("repo: could not get rate value from reponse")
	}

	return big.NewFloat(val.(float64)), nil
}

func (r *Repository) Name() string {
	return "coingecko"
}
