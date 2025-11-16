package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/model/dto"
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
		request    dto.CreatePayment
		statusCode int
		err        error
	)

	err = c.BindJSON(&request)
	if err != nil {
		statusCode = http.StatusBadRequest
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	payment, err := h.paymentUsecase.Create(c.Request.Context(), &request)
	if err != nil {
		statusCode = http.StatusNotFound
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusOK
	c.JSON(statusCode, model.ResponseSuccess(statusCode, payment))
}
