package application

import (
	"app.initial/src/orders/domain/entities"
	"app.initial/src/orders/domain/repositories"
)

type GetOrderUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewGetOrderUseCase(orderRepo repositories.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderUseCase) Execute(id uint) (*entities.Order, error) {
	return uc.orderRepo.FindByID(id)
}
