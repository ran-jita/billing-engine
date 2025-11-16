package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/usecase"
	"net/http"
)

type LoanHandler struct {
	loanUsecase *usecase.LoanUsecase
}

func NewLoanHandler(loanUsecase *usecase.LoanUsecase) *LoanHandler {
	return &LoanHandler{loanUsecase: loanUsecase}
}

func (h *LoanHandler) GetAll(c *gin.Context) {
	var (
		borrowerId string
		statusCode int
		err        error
	)

	borrowerId = c.Query("borrower_id")
	if borrowerId == "" {
		statusCode = http.StatusForbidden
		err = errors.New("borrower_id is required")
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	data, err := h.loanUsecase.GetAll(c.Request.Context(), borrowerId)
	if err != nil {
		statusCode = http.StatusBadRequest
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusOK
	c.JSON(statusCode, model.ResponseSuccess(statusCode, data))
}

func (h *LoanHandler) GetById(c *gin.Context) {
	var (
		loanId     string
		statusCode int
		err        error
	)

	loanId = c.Param("id")
	if loanId == "" {
		statusCode = http.StatusForbidden
		err = errors.New("loan id is required")
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	data, err := h.loanUsecase.GetById(c.Request.Context(), loanId)
	if err != nil {
		statusCode = http.StatusNotFound
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusOK
	c.JSON(statusCode, model.ResponseSuccess(statusCode, data))
}
