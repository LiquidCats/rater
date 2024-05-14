package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rater/internal/adapter/api/dto"
	"rater/internal/app/domain/types"
	"rater/internal/app/usecase"
)

type RateHandler struct {
	usecase *usecase.RateUsecase
}

func NewRateHandler(usecase *usecase.RateUsecase) *RateHandler {
	return &RateHandler{
		usecase: usecase,
	}
}

func (r *RateHandler) GetRate(ctx *gin.Context) {
	quote := ctx.Param("quote")
	base := ctx.Param("base")

	rate, err := r.usecase.GetRate(ctx, types.QuoteCurrency(quote), types.BaseCurrency(base))
	if nil != err {
		ctx.JSON(http.StatusBadRequest, dto.NewErrorResponse(err))

		return
	}

	ctx.JSON(http.StatusOK, dto.NewRateResponse(rate))
}
