package prometheus

import (
	"strconv"
	"time"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ProviderErrRate struct {
	*prometheus.CounterVec
}

func NewProviderErrRate(namespace string) *ProviderErrRate {
	counterVec := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "provider_err_rate",
	}, []string{"code", "provider"})

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
	histogramVec := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "response_time",
		Buckets:   []float64{0.1, 0.3, 0.5, 0.7, 1, 3, 5, 8, 13},
	}, []string{"route"})
	return &ResponseTime{histogramVec}
}

func (r *ResponseTime) Observe(route string, start time.Time) {
	r.WithLabelValues(route).Observe(
		float64(time.Since(start).Milliseconds()),
	)
}
