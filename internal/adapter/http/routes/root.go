package routes

import (
	"net/http"

	"github.com/LiquidCats/rater/internal/adapter/http/dto"
	"github.com/gin-gonic/gin"
)

type RootHandler struct {
}

func NewRootHandler() *RootHandler {
	return &RootHandler{}
}

func (r *RootHandler) Handle(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, dto.NewRootResponse("v1"))
}
