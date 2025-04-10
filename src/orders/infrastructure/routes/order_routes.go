package routes

import (
	"app.initial/src/core/container"
	"app.initial/src/orders/infrastructure/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(api *gin.RouterGroup, db *gorm.DB, rabbitMQURL string) {
	// Inicializar contenedor de dependencias
	c := container.NewOrdersContainer(db)

	// Inicializar consumidor de eventos
	if err := c.ConfigureEventConsumer(rabbitMQURL); err != nil {
		panic(err)
	}

	// Inicializar controladores con casos de uso inyectados
	createOrderController := controllers.NewCreateOrderController(
		c.GetCreateOrderUseCase(rabbitMQURL),
	)
	getOrderController := controllers.NewGetOrderController(
		c.GetGetOrderUseCase(),
	)
	getCustomerOrdersController := controllers.NewGetCustomerOrdersController(
		c.GetGetCustomerOrdersUseCase(),
	)

	// Configurar rutas
	orders := api.Group("/orders")
	{
		orders.POST("", createOrderController.Handle)
		orders.GET("/:id", getOrderController.Handle)
		orders.GET("/customer/:customerID", getCustomerOrdersController.Handle)
	}
}
