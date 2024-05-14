package coinapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"rater/internal/app/domain/types"
	"time"
)

type Repository struct {
	url      string
	apiToken string
}

type coinApiResponse struct {
	Time         time.Time `json:"time"`
	AssetIdBase  string    `json:"asset_id_base"`
	AssetIdQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}

func NewRepository() *Repository {
	url := os.Getenv("RATER_COINAPI_URL")
	apiToken := os.Getenv("RATER_COINAPI_SECRET")

	return &Repository{
		url:      url,
		apiToken: apiToken,
	}
}

func (a *Repository) Get(ctx context.Context, quote types.QuoteCurrency, base types.BaseCurrency) (*big.Float, error) {
	url := fmt.Sprintf(
		"%s/%s/%s",
		a.url,
		base.Upper(),
		quote.Upper(),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("repo: could not create request: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CoinAPI-Key", a.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo: error making api request: %s\n", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("repo: error making api request: %s\n", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("repo: could not read response body: %s\n", err)
	}

	var resp coinApiResponse

	if err = json.Unmarshal(resBody, &resp); nil != err {
		return nil, fmt.Errorf("repo: could not parse response: %s\n", err)
	}

	return big.NewFloat(resp.Rate), nil
}

func (a *Repository) Name() string {
	return "coinapi"
}
