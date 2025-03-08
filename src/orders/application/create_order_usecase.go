package application

import (
	"app.initial/src/orders/domain/entities"
	"app.initial/src/orders/domain/repositories"
	"github.com/google/uuid"
	"time"
)

type CreateOrderUseCase struct {
	orderRepo      repositories.OrderRepository
	eventPublisher repositories.EventPublisher
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	eventPublisher repositories.EventPublisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (uc *CreateOrderUseCase) Execute(customerID uint, items []entities.OrderItem) (*entities.Order, error) {
	// Calculate total amount
	var totalAmount float64
	for _, item := range items {
		totalAmount += item.TotalPrice
	}

	// Create new order
	order := &entities.Order{
		CustomerID:  customerID,
		Items:       items,
		TotalAmount: totalAmount,
		Status:      entities.OrderStatusCreated,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository
	if err := uc.orderRepo.Save(order); err != nil {
		return nil, err
	}

	// Publish order.created event
	event := entities.Event{
		ID:        uuid.New().String(),
		Type:      "order.created",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"order_id":     order.ID,
			"customer_id":  order.CustomerID,
			"total_amount": order.TotalAmount,
		},
	}

	if err := uc.eventPublisher.PublishEvent(event); err != nil {
		// Log error but continue
		// In a production environment, we might want to implement a retry mechanism
		// or store failed events for later processing
	}

	return order, nil
}
