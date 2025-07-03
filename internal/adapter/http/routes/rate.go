package routes

import (
	"net/http"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type RateHandler struct {
	usecase *usecase.RateUsecase
}

func NewRateHandler(usecase *usecase.RateUsecase) *RateHandler {
	return &RateHandler{
		usecase: usecase,
	}
}

func (r *RateHandler) Handle(ctx *gin.Context) {
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
