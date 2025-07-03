package metrics

import "github.com/LiquidCats/rater/internal/app/domain/entity"

type ProviderErrRateMetric interface {
	Inc(code int, provider entity.ProviderName)
}
