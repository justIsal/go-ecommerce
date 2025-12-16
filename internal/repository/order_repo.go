package repository

import (
	"errors"
	"go-ecommerce/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository interface {
	CreateOrder(order *domain.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *domain.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		
		if err := tx.Create(order).Error; err != nil {
			return err 
		}

		for _, item := range order.Items {
			var product domain.Product
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&product, item.ProductID).Error; err != nil {
				return errors.New("product not found")
			}

			if product.Stock < item.Quantity {
				return errors.New("insufficient stock for product: " + product.Name)
			}

			product.Stock -= item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				return err
			}
		}

		return nil
	})
}