package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"net/http"
)

type PingHandler struct{}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (c *PingHandler) Ping(ctx *gin.Context) {
	statusCode := http.StatusOK

	ctx.JSON(statusCode, model.ResponseSuccess(
		statusCode,
		"ping",
	))
}
