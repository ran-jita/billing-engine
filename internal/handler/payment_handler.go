package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/usecase"
	"net/http"
)

type PaymentHandler struct {
	paymentUsecase *usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase *usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{paymentUsecase: paymentUsecase}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var (
		payment    model.Payment
		statusCode int
		err        error
	)

	err = c.BindJSON(&payment)
	if err != nil {
		statusCode = http.StatusBadRequest
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	err = h.paymentUsecase.Create(c.Request.Context(), &payment)
	if err != nil {
		statusCode = http.StatusNotFound
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusOK
	c.JSON(statusCode, model.ResponseSuccess(statusCode, payment))
}
