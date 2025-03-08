package repositories

import (
	"app.initial/src/orders/domain/entities"
	"gorm.io/gorm"
)

type MySQLOrderRepository struct {
	db *gorm.DB
}

func NewMySQLOrderRepository(db *gorm.DB) *MySQLOrderRepository {
	return &MySQLOrderRepository{db: db}
}

func (r *MySQLOrderRepository) Save(order *entities.Order) error {
	return r.db.Save(order).Error
}

func (r *MySQLOrderRepository) FindByID(id uint) (*entities.Order, error) {
	var order entities.Order
	result := r.db.Preload("Items").First(&order, id)
	return &order, result.Error
}

func (r *MySQLOrderRepository) FindAll() ([]entities.Order, error) {
	var orders []entities.Order
	result := r.db.Preload("Items").Find(&orders)
	return orders, result.Error
}

func (r *MySQLOrderRepository) FindByCustomerID(customerID uint) ([]entities.Order, error) {
	var orders []entities.Order
	result := r.db.Preload("Items").Where("customer_id = ?", customerID).Find(&orders)
	return orders, result.Error
}

func (r *MySQLOrderRepository) UpdateStatus(id uint, status entities.OrderStatus) error {
	return r.db.Model(&entities.Order{}).Where("id = ?", id).Update("status", status).Error
}
