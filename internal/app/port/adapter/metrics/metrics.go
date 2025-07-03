package metrics

import (
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
)

type ProviderErrRateMetric interface {
	Inc(code int, provider entity.ProviderName)
}

type ResponseTimeMetric interface {
	Observe(route string, start time.Time)
}
