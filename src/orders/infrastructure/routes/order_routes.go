package routes

import (
	"app.initial/src/orders/application"
	"app.initial/src/orders/infrastructure/controllers"
	"app.initial/src/orders/infrastructure/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupOrderRoutes(api *gin.RouterGroup, db *gorm.DB, rabbitMQURL string) {
	// Initialize repositories
	orderRepo := repositories.NewMySQLOrderRepository(db)
	eventPublisher, err := repositories.NewRabbitMQEventPublisher(rabbitMQURL, "orders_exchange")
	if err != nil {
		panic(err)
	}

	// Initialize use cases
	createOrderUseCase := application.NewCreateOrderUseCase(orderRepo, eventPublisher)
	getOrderUseCase := application.NewGetOrderUseCase(orderRepo)
	getCustomerOrdersUseCase := application.NewGetCustomerOrdersUseCase(orderRepo)
	processPaymentUseCase := application.NewProcessPaymentUseCase(orderRepo, eventPublisher)

	// Initialize event consumer
	eventConsumer, err := repositories.NewRabbitMQEventConsumer(rabbitMQURL, processPaymentUseCase)
	if err != nil {
		panic(err)
	}

	// Start consuming payment events
	err = eventConsumer.StartConsumingPaymentEvents(
		"payment_events_queue",
		"payments_exchange",
		"payment.#",
	)
	if err != nil {
		panic(err)
	}

	// Initialize controllers
	createOrderController := controllers.NewCreateOrderController(createOrderUseCase)
	getOrderController := controllers.NewGetOrderController(getOrderUseCase)
	getCustomerOrdersController := controllers.NewGetCustomerOrdersController(getCustomerOrdersUseCase)

	// Setup routes
	orders := api.Group("/orders")
	{
		orders.POST("", createOrderController.Handle)
		orders.GET("/:id", getOrderController.Handle)
		orders.GET("/customer/:customerID", getCustomerOrdersController.Handle)
	}
}
