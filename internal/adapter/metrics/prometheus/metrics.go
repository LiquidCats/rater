package prometheus

import (
	"strconv"
	"time"

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

type ResponseTime struct {
	*prometheus.HistogramVec
}

func NewResponseTime(namespace string) *ResponseTime {
	histogramVec := prometheus.V2.NewHistogramVec(prometheus.HistogramVecOpts{
		HistogramOpts: prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "response_time",
			Buckets:   []float64{0.1, 0.3, 0.5, 0.7, 1, 3, 5, 8, 13},
		},
		VariableLabels: prometheus.ConstrainedLabels{
			{Name: "route"},
		},
	})
	return &ResponseTime{histogramVec}
}

func (r *ResponseTime) Observe(route string, start time.Time) {
	r.WithLabelValues(route).Observe(
		float64(time.Since(start).Milliseconds()),
	)
}
