package routes

import (
	"net/http"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/usecase"

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

	pair, err := pairStr.ToPair()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewErrorResponse(err))
		return
	}

	rate, err := r.usecase.GetRate(ctx, pair)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.NewErrorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, dto.NewRateResponse(rate))
}
