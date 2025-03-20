package container

import (
	"log"

	"app.initial/src/orders/application"
	"app.initial/src/orders/domain/repositories"
	infraRepo "app.initial/src/orders/infrastructure/repositories"
	"gorm.io/gorm"
)

// OrdersContainer centraliza la creación y gestión de dependencias para órdenes
type OrdersContainer struct {
	db                       *gorm.DB
	orderRepo                repositories.OrderRepository
	eventPublisher           repositories.EventPublisher
	createOrderUseCase       *application.CreateOrderUseCase
	getOrderUseCase          *application.GetOrderUseCase
	getCustomerOrdersUseCase *application.GetCustomerOrdersUseCase
	processPaymentUseCase    *application.ProcessPaymentUseCase
}

// NewOrdersContainer crea un nuevo contenedor
func NewOrdersContainer(db *gorm.DB) *OrdersContainer {
	return &OrdersContainer{
		db: db,
	}
}

// GetOrderRepository devuelve el repositorio de órdenes
func (c *OrdersContainer) GetOrderRepository() repositories.OrderRepository {
	if c.orderRepo == nil {
		c.orderRepo = infraRepo.NewMySQLOrderRepository(c.db)
	}
	return c.orderRepo
}

// GetEventPublisher devuelve el publicador de eventos
func (c *OrdersContainer) GetEventPublisher(rabbitMQURL string) repositories.EventPublisher {
	if c.eventPublisher == nil {
		publisher, err := infraRepo.NewRabbitMQEventPublisher(rabbitMQURL, "orders_exchange")
		if err != nil {
			log.Fatalf("Error creating event publisher: %v", err)
		}
		c.eventPublisher = publisher
	}
	return c.eventPublisher
}

// GetCreateOrderUseCase devuelve el caso de uso de crear órdenes
func (c *OrdersContainer) GetCreateOrderUseCase(rabbitMQURL string) *application.CreateOrderUseCase {
	if c.createOrderUseCase == nil {
		c.createOrderUseCase = application.NewCreateOrderUseCase(
			c.GetOrderRepository(),
			c.GetEventPublisher(rabbitMQURL),
		)
	}
	return c.createOrderUseCase
}

// GetGetOrderUseCase devuelve el caso de uso de obtener órdenes
func (c *OrdersContainer) GetGetOrderUseCase() *application.GetOrderUseCase {
	if c.getOrderUseCase == nil {
		c.getOrderUseCase = application.NewGetOrderUseCase(
			c.GetOrderRepository(),
		)
	}
	return c.getOrderUseCase
}

// GetGetCustomerOrdersUseCase devuelve el caso de uso de obtener órdenes por cliente
func (c *OrdersContainer) GetGetCustomerOrdersUseCase() *application.GetCustomerOrdersUseCase {
	if c.getCustomerOrdersUseCase == nil {
		c.getCustomerOrdersUseCase = application.NewGetCustomerOrdersUseCase(
			c.GetOrderRepository(),
		)
	}
	return c.getCustomerOrdersUseCase
}

// GetProcessPaymentUseCase devuelve el caso de uso de procesar pagos
func (c *OrdersContainer) GetProcessPaymentUseCase(rabbitMQURL string) *application.ProcessPaymentUseCase {
	if c.processPaymentUseCase == nil {
		c.processPaymentUseCase = application.NewProcessPaymentUseCase(
			c.GetOrderRepository(),
			c.GetEventPublisher(rabbitMQURL),
		)
	}
	return c.processPaymentUseCase
}

// ConfigureEventConsumer configura el consumidor de eventos
func (c *OrdersContainer) ConfigureEventConsumer(rabbitMQURL string) error {
	consumer, err := infraRepo.NewRabbitMQEventConsumer(
		rabbitMQURL,
		c.GetProcessPaymentUseCase(rabbitMQURL),
	)
	if err != nil {
		return err
	}

	return consumer.StartConsumingPaymentEvents(
		"payment_events_queue",
		"payments_exchange",
		"payment.#",
	)
}
