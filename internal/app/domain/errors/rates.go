package errors

import "github.com/pkg/errors"

var (
	ErrRateNotAvailable = errors.New("exchange rate is not available right now")
)
