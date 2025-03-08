package application

import (
	"app.initial/src/orders/domain/entities"
	"app.initial/src/orders/domain/repositories"
)

type GetCustomerOrdersUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewGetCustomerOrdersUseCase(orderRepo repositories.OrderRepository) *GetCustomerOrdersUseCase {
	return &GetCustomerOrdersUseCase{orderRepo: orderRepo}
}

func (uc *GetCustomerOrdersUseCase) Execute(customerID uint) ([]entities.Order, error) {
	return uc.orderRepo.FindByCustomerID(customerID)
}
