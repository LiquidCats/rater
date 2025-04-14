package middlware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/gin-gonic/gin"
)

type PairValidationMiddleware struct {
	allowed []entity.CurrencyPairString
}

func NewPairValidation(allowed []entity.CurrencyPairString) *PairValidationMiddleware {
	return &PairValidationMiddleware{
		allowed: allowed,
	}
}

func (m *PairValidationMiddleware) Handle(ctx *gin.Context) {
	pair := strings.TrimSpace(ctx.Param("pair"))

	if len(m.allowed) == 0 {
		m.setValidationError(ctx, "no allowed pairs")
		return
	}

	if pair == "" {
		m.setValidationError(ctx, "empty pair")
		ctx.Abort()
		return
	}
	if len(pair) < 7 { // nolint:mnd
		m.setValidationError(ctx, "should be 7 characters long at least")
		return
	}

	if pair != strings.ToUpper(pair) {
		m.setValidationError(ctx, "should be uppercase")
		return
	}

	if !slices.Contains(m.allowed, entity.CurrencyPairString(pair)) {
		m.setValidationError(ctx, "not allowed")
		return
	}

	ctx.Next()
}

func (m *PairValidationMiddleware) setValidationError(ctx *gin.Context, text ...string) {
	ctx.AbortWithStatusJSON(
		http.StatusUnprocessableEntity,
		dto.NewValidationErrorResponse(
			"pair",
			text...,
		),
	)
}
