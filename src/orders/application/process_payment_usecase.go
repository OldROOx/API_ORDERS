package application

import (
	"app.initial/src/orders/domain/entities"
	"app.initial/src/orders/domain/repositories"
	"github.com/google/uuid"
	"time"
)

type ProcessPaymentUseCase struct {
	orderRepo      repositories.OrderRepository
	eventPublisher repositories.EventPublisher
}

func NewProcessPaymentUseCase(
	orderRepo repositories.OrderRepository,
	eventPublisher repositories.EventPublisher,
) *ProcessPaymentUseCase {
	return &ProcessPaymentUseCase{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (uc *ProcessPaymentUseCase) PaymentCompleted(orderID uint, paymentID string) error {

	if err := uc.orderRepo.UpdateStatus(orderID, entities.OrderStatusPaid); err != nil {
		return err
	}

	event := entities.Event{
		ID:        uuid.New().String(),
		Type:      "order.paid",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"order_id":   orderID,
			"payment_id": paymentID,
		},
	}

	return uc.eventPublisher.PublishEvent(event)
}

func (uc *ProcessPaymentUseCase) PaymentFailed(orderID uint, reason string) error {

	if err := uc.orderRepo.UpdateStatus(orderID, entities.OrderStatusCanceled); err != nil {
		return err
	}

	event := entities.Event{
		ID:        uuid.New().String(),
		Type:      "order.payment_failed",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"order_id": orderID,
			"reason":   reason,
		},
	}

	return uc.eventPublisher.PublishEvent(event)
}
