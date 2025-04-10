package repositories

import (
	"app.initial/src/orders/domain/entities"
)

type EventPublisher interface {
	PublishEvent(event entities.Event) error
	Close() error
}
