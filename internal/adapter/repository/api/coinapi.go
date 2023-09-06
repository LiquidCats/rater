package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

type CoinApiRepository struct {
	url      string
	apiToken string
}

type coinApiResponse struct {
	Time         time.Time `json:"time"`
	AssetIdBase  string    `json:"asset_id_base"`
	AssetIdQuote string    `json:"asset_id_quote"`
	Rate         float64   `json:"rate"`
}

func NewCoinApiRepository(url, apiToken string) *CoinApiRepository {
	return &CoinApiRepository{
		url:      url,
		apiToken: apiToken,
	}
}

func (a *CoinApiRepository) Get(ctx context.Context, quote, base string) (*big.Float, error) {
	url := fmt.Sprintf("%s/%s/%s", a.url, strings.ToUpper(base), strings.ToUpper(quote))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("repo: could not create request: %s\n", err)
	}

	req.WithContext(ctx)
	req.Header.Set("X-CoinAPI-Key", a.apiToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("repo: error making api request: %s\n", err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("repo: could not read response body: %s\n", err)
	}

	var resp coinApiResponse

	if err = json.Unmarshal(resBody, &resp); nil != err {
		return nil, fmt.Errorf("repo: could not parse response: %s\n", err)
	}

	if err != nil {
		return nil, fmt.Errorf("repo: could not parse response: %s\n", err)
	}

	return big.NewFloat(resp.Rate), nil
}

func (a *CoinApiRepository) Name() string {
	return "coinapi"
}
