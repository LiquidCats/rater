package middlware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rater/internal/adapter/api/dto"
	"rater/internal/app/domain/types"
	"rater/pkg/utils/sliceutils"
	"strings"
)

type CurrencyParamMiddleware struct {
	allowed []types.CurrencyNameString
	key     string
}

func NewCurrencyParamMiddleware(key string, allowed []types.CurrencyNameString) *CurrencyParamMiddleware {
	return &CurrencyParamMiddleware{
		allowed: allowed,
		key:     key,
	}
}

func (i *CurrencyParamMiddleware) Handle(ctx *gin.Context) {
	base := ctx.Param(i.key)
	if 3 != len(base) {
		ctx.JSON(http.StatusUnprocessableEntity, dto.NewFailResponse(gin.H{
			i.key: "should be 3 characters long",
		}))
		ctx.Abort()
		return
	}

	if base != strings.ToUpper(base) {
		ctx.JSON(http.StatusUnprocessableEntity, dto.NewFailResponse(gin.H{
			i.key: "should be uppercase",
		}))
		ctx.Abort()
		return
	}

	if !sliceutils.Contains(i.allowed, types.CurrencyNameString(base)) {
		ctx.JSON(http.StatusUnprocessableEntity, dto.NewFailResponse(gin.H{
			i.key: "not allowed",
		}))
		ctx.Abort()
		return
	}

	ctx.Next()
}
