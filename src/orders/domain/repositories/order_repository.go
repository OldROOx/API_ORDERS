package repositories

import (
	"app.initial/src/orders/domain/entities"
)

type OrderRepository interface {
	Save(order *entities.Order) error
	FindByID(id uint) (*entities.Order, error)
	FindAll() ([]entities.Order, error)
	FindByCustomerID(customerID uint) ([]entities.Order, error)
	UpdateStatus(id uint, status entities.OrderStatus) error
}
