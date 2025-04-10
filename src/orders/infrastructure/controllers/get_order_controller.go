package controllers

import (
	"app.initial/src/orders/application"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetOrderController struct {
	getOrderUseCase *application.GetOrderUseCase
}

func NewGetOrderController(getOrderUseCase *application.GetOrderUseCase) *GetOrderController {
	return &GetOrderController{getOrderUseCase: getOrderUseCase}
}

func (c *GetOrderController) Handle(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	order, err := c.getOrderUseCase.Execute(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	ctx.JSON(http.StatusOK, order)
}
