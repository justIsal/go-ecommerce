package service

import (
	"errors"
	"go-ecommerce/internal/domain"
	"go-ecommerce/internal/repository"
)

type CartItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

type OrderService interface {
	CreateOrder(userID uint, items []CartItemInput) (*domain.Order, error)
}

type orderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &orderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *orderService) CreateOrder(userID uint, inputItems []CartItemInput) (*domain.Order, error) {
	var orderItems []domain.OrderItem
	var totalPrice float64

	for _, item := range inputItems {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return nil, errors.New("product not found")
		}

		if product.Stock < item.Quantity {
			return nil, errors.New("insufficient stock for " + product.Name)
		}

		itemPrice := product.Price
		subTotal := itemPrice * float64(item.Quantity)
		totalPrice += subTotal

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: product.ID,
			Quantity:  item.Quantity,
			Price:     itemPrice, 
		})
	}

	order := domain.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
		Items:      orderItems, 
	}

	if err := s.orderRepo.CreateOrder(&order); err != nil {
		return nil, err
	}

	return &order, nil
}