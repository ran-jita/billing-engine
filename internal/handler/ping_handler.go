package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"net/http"
)

type PingHandler interface {
	Ping(ctx *gin.Context)
}

type pingHandler struct{}

func NewPingHandler() *pingHandler {
	return &pingHandler{}
}

func (c *pingHandler) Ping(ctx *gin.Context) {
	statusCode := http.StatusOK

	ctx.JSON(statusCode, model.ResponseSuccess(
		statusCode,
		"ping",
	))
}
