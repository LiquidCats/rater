package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"rater/internal/adapter/api/dto"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (r *RootHandler) GetRoot(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.NewRootResponse("v1"))
}
