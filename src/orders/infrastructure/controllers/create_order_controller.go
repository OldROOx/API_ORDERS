package controllers

import (
	"app.initial/src/orders/application"
	"app.initial/src/orders/domain/entities"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateOrderRequest struct {
	CustomerID uint                 `json:"customer_id" binding:"required"`
	Items      []entities.OrderItem `json:"items" binding:"required"`
}

type CreateOrderController struct {
	createOrderUseCase *application.CreateOrderUseCase
}

func NewCreateOrderController(createOrderUseCase *application.CreateOrderUseCase) *CreateOrderController {
	return &CreateOrderController{createOrderUseCase: createOrderUseCase}
}

func (c *CreateOrderController) Handle(ctx *gin.Context) {
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := c.createOrderUseCase.Execute(req.CustomerID, req.Items)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}
