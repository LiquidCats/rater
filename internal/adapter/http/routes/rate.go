package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/port/adapter/metrics"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	usecase *usecase.RateUseCase
	metrics Metrics
}

type Metrics struct {
	ResponseTime metrics.ResponseTimeMetric
}

func NewRateHandler(usecase *usecase.RateUseCase, metrics Metrics) *RateHandler {
	return &RateHandler{
		usecase: usecase,
		metrics: metrics,
	}
}

func (r *RateHandler) Handle(ctx *gin.Context) {
	start := time.Now()
	defer r.metrics.ResponseTime.Observe(ctx.Request.URL.Path, start)

	symbol := entity.Symbol(ctx.Param("pair"))

	date := time.Now()
	if d := strings.TrimSpace(ctx.Query("date")); d != "" {
		date, _ = time.Parse(entity.DefaultFormat, d)
	}

	logger := zerolog.Ctx(ctx.Request.Context()).
		With().
		Any("symbol", symbol).
		Logger()

	rate, err := r.usecase.GetRate(ctx, symbol, date)
	if err != nil {
		logger.Error().
			Any("err", eris.ToJSON(err, true)).
			Msg("invalid rate")

		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewRateResponse(rate))
}
