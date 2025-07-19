package middlware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/port/adapter/repository"
	"github.com/gin-gonic/gin"
)

type PairValidationMiddleware struct {
	db repository.PairDatabase
}

func NewPairValidation(db repository.PairDatabase) *PairValidationMiddleware {
	return &PairValidationMiddleware{
		db: db,
	}
}

func (m *PairValidationMiddleware) Handle(ctx *gin.Context) {
	pair := strings.TrimSpace(ctx.Param("pair"))

	count, err := m.db.CountPairs(ctx.Request.Context())
	if err != nil {
		setValidationError(ctx, fmt.Sprintf("pair validation failed (%s)", err.Error()))
		return
	}

	if count == 0 {
		setValidationError(ctx, "no allowed pairs")
		return
	}

	if pair == "" {
		setValidationError(ctx, "empty pair")
		ctx.Abort()
		return
	}
	if len(pair) < 7 { // nolint:mnd
		setValidationError(ctx, "should be 7 characters long at least")
		return
	}

	if pair != strings.ToUpper(pair) {
		setValidationError(ctx, "should be uppercase")
		return
	}

	exists, err := m.db.HasPair(ctx, strings.ToUpper(pair))
	if err != nil {
		setValidationError(ctx, fmt.Sprintf("pair validation failed (%s)", err.Error()))
		return
	}

	if !exists {
		setValidationError(ctx, "not allowed")
		return
	}

	ctx.Next()
}

func setValidationError(ctx *gin.Context, text ...string) {
	ctx.AbortWithStatusJSON(
		http.StatusUnprocessableEntity,
		dto.NewValidationErrorResponse(
			"pair",
			text...,
		),
	)
}
