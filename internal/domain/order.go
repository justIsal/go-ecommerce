package domain

import (
	"time"
)

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null" json:"user_id"` 
	TotalPrice float64        `gorm:"type:decimal(15,2)" json:"total_price"`
	Status     string         `gorm:"type:varchar(20);default:'pending'" json:"status"` 
	Items      []OrderItem    `gorm:"foreignKey:OrderID" json:"items"` 
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Product   Product `json:"product"` 
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(15,2)" json:"price"` 
}