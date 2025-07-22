package middlware

import (
	"net/http"
	"strings"
	"time"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/LiquidCats/rater/internal/app/utils/timeutils"
	"github.com/gin-gonic/gin"
)

type DateValidationMiddleware struct {
}

func NewDateValidation() *DateValidationMiddleware {
	return &DateValidationMiddleware{}
}

func (m *DateValidationMiddleware) Handle(ctx *gin.Context) {
	date := strings.TrimSpace(ctx.Query("date"))
	if date == "" {
		ctx.Next()
		return
	}

	ts, err := time.Parse(entity.DefaultFormat, date)
	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			dto.NewValidationErrorResponse(
				"date",
				"datetime must be in format: "+entity.DefaultFormat,
				err.Error(),
			),
		)
		return
	}

	now := time.Now()
	roundedNow := timeutils.RoundToNearest(
		time.Date(
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
			now.Nanosecond(),
			time.UTC,
		),
		timeutils.FiveMinuteBucket,
	)
	rounded := timeutils.RoundToNearest(ts, timeutils.FiveMinuteBucket)
	if rounded.Compare(roundedNow) > 0 {
		ctx.AbortWithStatusJSON(
			http.StatusUnprocessableEntity,
			dto.NewValidationErrorResponse(
				"date",
				"datetime must be in the past",
			),
		)
	}

	ctx.Next()
}
