package dto

import (
	"github.com/gin-gonic/gin"
	"rater/internal/app/domain/entity"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
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
		"message": err,
	}
}

func NewRootResponse(version string) gin.H {
	return gin.H{"version": version}
}

func NewRateResponse(rate *entity.Rate) gin.H {
	return NewSuccessResponse(gin.H{
		"quote": rate.Quote,
		"base":  rate.Base,
		"price": rate.Price,
	})
}
