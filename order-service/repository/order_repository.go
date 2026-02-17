package repository

import (
	"context"

	"gorm.io/gorm"

	"online-store-microservice/order-service/models"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
	GetByUserID(ctx context.Context, userID string) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
