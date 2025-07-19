package errors

import (
	"fmt"

	"github.com/rotisserie/eris"
)

var (
	ErrRateNotAvailable = eris.New("exchange rate is not available right now")
	ErrNoHistoricalRate = eris.New("historical rate is not available for this time")
)

type ProviderRequestFailedError struct {
	StatusCode int
	Body       string
}

func (e ProviderRequestFailedError) Error() string {
	return fmt.Sprintf("provider request failed with code %d", e.StatusCode)
}
