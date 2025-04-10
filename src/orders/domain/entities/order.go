package entities

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusCreated  OrderStatus = "CREATED"
	OrderStatusPending  OrderStatus = "PENDING"
	OrderStatusPaid     OrderStatus = "PAID"
	OrderStatusShipped  OrderStatus = "SHIPPED"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	OrderID    uint    `json:"order_id" gorm:"not null"`
	ProductID  uint    `json:"product_id" gorm:"not null"`
	Name       string  `json:"name" gorm:"size:100;not null"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	UnitPrice  float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	TotalPrice float64 `json:"total_price" gorm:"type:decimal(10,2);not null"`
}

type Order struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	CustomerID  uint        `json:"customer_id" gorm:"not null"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	TotalAmount float64     `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Status      OrderStatus `json:"status" gorm:"size:20;not null"`
	PaymentID   string      `json:"payment_id,omitempty" gorm:"size:100"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}
