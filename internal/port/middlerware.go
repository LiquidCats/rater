package port

import "github.com/gin-gonic/gin"

type Middleware interface {
	Handle(ctx *gin.Context)
}
