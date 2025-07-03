package prometheus

import (
	"strconv"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/prometheus/client_golang/prometheus"
)

type ProviderErrRate struct {
	*prometheus.CounterVec
}

func NewProviderErrRate(namespace string) *ProviderErrRate {
	counterVec := prometheus.V2.NewCounterVec(prometheus.CounterVecOpts{
		CounterOpts: prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "provider_err_rate",
		},
		VariableLabels: prometheus.ConstrainedLabels{
			{Name: "code"},
			{Name: "provider"},
		},
	})

	prometheus.MustRegister(counterVec)

	return &ProviderErrRate{
		CounterVec: counterVec,
	}
}

func (p *ProviderErrRate) Inc(code int, provider entity.ProviderName) {
	p.WithLabelValues(strconv.Itoa(code), string(provider)).Inc()
}
