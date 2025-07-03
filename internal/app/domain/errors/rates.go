package errors

import (
	"fmt"

	"github.com/rotisserie/eris"
)

var (
	ErrRateNotAvailable = eris.New("exchange rate is not available right now")
)

type ErrProviderRequestFailed struct {
	StatusCode int
	Body       string
}

func (e ErrProviderRequestFailed) Error() string {
	return fmt.Sprintf("provider request failed with code %d", e.StatusCode)
}
