package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ran-jita/billing-engine/internal/model"
	"github.com/ran-jita/billing-engine/internal/model/dto"
	"github.com/ran-jita/billing-engine/internal/usecase"
	"net/http"
	"time"
)

type BorrowerHandler struct {
	borrowerUsecase *usecase.BorrowerUsecase
}

func NewBorrowerHandler(borrowerUsecase *usecase.BorrowerUsecase) *BorrowerHandler {
	return &BorrowerHandler{borrowerUsecase: borrowerUsecase}
}

func (h *BorrowerHandler) GetById(c *gin.Context) {
	var (
		borrowerId string
		statusCode int
		err        error
	)

	borrowerId = c.Param("id")
	if borrowerId == "" {
		statusCode = http.StatusForbidden
		err = errors.New("loan id is required")
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	data, err := h.borrowerUsecase.GetById(c.Request.Context(), borrowerId)
	if err != nil {
		statusCode = http.StatusNotFound
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusOK
	c.JSON(statusCode, model.ResponseSuccess(statusCode, data))
}

func (h *BorrowerHandler) UpdateStatusDelinquent(c *gin.Context) {
	var (
		request           dto.UpdateDelinquent
		statusCode        int
		err               error
		parsedProcessDate time.Time
	)

	err = c.BindJSON(&request)
	if err != nil {
		statusCode = http.StatusBadRequest
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	parsedProcessDate = time.Now()
	if request.ProcessDate != "" {
		parsedProcessDate, err = time.Parse("2006-01-02", request.ProcessDate)
	}

	err = h.borrowerUsecase.UpdateDelinquent(c.Request.Context(), parsedProcessDate)
	if err != nil {
		statusCode = http.StatusNotFound
		c.JSON(statusCode, model.ResponseError(statusCode, err))
		return
	}

	statusCode = http.StatusNoContent
	c.JSON(statusCode, model.ResponseSuccess(statusCode, nil))
}
