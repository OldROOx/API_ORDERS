package routes

import (
	"app.initial/src/orders/application"
	"app.initial/src/orders/infrastructure/controllers"
	"app.initial/src/orders/infrastructure/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(api *gin.RouterGroup, db *gorm.DB, rabbitMQURL string) {

	orderRepo := repositories.NewMySQLOrderRepository(db)
	eventPublisher, err := repositories.NewRabbitMQEventPublisher(rabbitMQURL, "orders_exchange")
	if err != nil {
		panic(err)
	}

	createOrderUseCase := application.NewCreateOrderUseCase(orderRepo, eventPublisher)
	getOrderUseCase := application.NewGetOrderUseCase(orderRepo)
	getCustomerOrdersUseCase := application.NewGetCustomerOrdersUseCase(orderRepo)
	processPaymentUseCase := application.NewProcessPaymentUseCase(orderRepo, eventPublisher)

	eventConsumer, err := repositories.NewRabbitMQEventConsumer(rabbitMQURL, processPaymentUseCase)
	if err != nil {
		panic(err)
	}

	err = eventConsumer.StartConsumingPaymentEvents(
		"payment_events_queue",
		"payments_exchange",
		"payment.#",
	)
	if err != nil {
		panic(err)
	}

	createOrderController := controllers.NewCreateOrderController(createOrderUseCase)
	getOrderController := controllers.NewGetOrderController(getOrderUseCase)
	getCustomerOrdersController := controllers.NewGetCustomerOrdersController(getCustomerOrdersUseCase)

	orders := api.Group("/orders")
	{
		orders.POST("", createOrderController.Handle)
		orders.GET("/:id", getOrderController.Handle)
		orders.GET("/customer/:customerID", getCustomerOrdersController.Handle)
	}
}
