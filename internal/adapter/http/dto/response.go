package dto

import (
	"github.com/LiquidCats/rater/internal/app/domain/entity"
	"github.com/gin-gonic/gin"
	"github.com/rotisserie/eris"
)

type ResponseStatus string

const (
	StatusSuccess ResponseStatus = "success"
	StatusFail    ResponseStatus = "fail"
	StatusError   ResponseStatus = "error"
)

func NewSuccessResponse(data gin.H) gin.H {
	return gin.H{
		"status": StatusSuccess,
		"data":   data,
	}
}

func NewFailResponse(data interface{}) gin.H {
	return gin.H{
		"status": StatusFail,
		"data":   data,
	}
}

func NewErrorResponse(err error) gin.H {
	return gin.H{
		"status":  StatusError,
		"message": eris.ToJSON(err, false),
	}
}

func NewRootResponse(version string) gin.H {
	return gin.H{"version": version}
}

func NewRateResponse(rate *entity.Rate) gin.H {
	return NewSuccessResponse(gin.H{
		"pair":  rate.Pair.Symbol.ToUpper().String(),
		"price": rate.Price.String(),
	})
}

func NewValidationErrorResponse(key string, texts ...string) gin.H {
	return NewFailResponse(gin.H{
		key: texts,
	})
}
