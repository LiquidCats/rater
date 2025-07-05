package routes

import (
	"net/http"
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
	usecase *usecase.RateUsecase
	metrics Metrics
}

type Metrics struct {
	ResponseTime metrics.ResponseTimeMetric
}

func NewRateHandler(usecase *usecase.RateUsecase, metrics Metrics) *RateHandler {
	return &RateHandler{
		usecase: usecase,
		metrics: metrics,
	}
}

func (r *RateHandler) Handle(ctx *gin.Context) {
	start := time.Now()
	defer r.metrics.ResponseTime.Observe(ctx.Request.URL.Path, start)

	pairStr := entity.CurrencyPairString(ctx.Param("pair"))

	logger := zerolog.Ctx(ctx.Request.Context()).With().Any("pair", pairStr).Logger()

	pair, err := pairStr.ToPair()
	if err != nil {
		logger.Error().Any("err", eris.ToJSON(err, true)).Msg("invalid pair")

		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	rate, err := r.usecase.GetRate(ctx, pair)
	if err != nil {
		logger.Error().Any("err", eris.ToJSON(err, true)).Msg("invalid rate")

		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, dto.NewRateResponse(rate))
}
