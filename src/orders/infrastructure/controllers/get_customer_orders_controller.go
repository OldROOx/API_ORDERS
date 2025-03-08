package controllers

import (
	"app.initial/src/orders/application"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GetCustomerOrdersController struct {
	getCustomerOrdersUseCase *application.GetCustomerOrdersUseCase
}

func NewGetCustomerOrdersController(getCustomerOrdersUseCase *application.GetCustomerOrdersUseCase) *GetCustomerOrdersController {
	return &GetCustomerOrdersController{getCustomerOrdersUseCase: getCustomerOrdersUseCase}
}

func (c *GetCustomerOrdersController) Handle(ctx *gin.Context) {
	customerID, err := strconv.ParseUint(ctx.Param("customerID"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	orders, err := c.getCustomerOrdersUseCase.Execute(uint(customerID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orders)
}
