package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"online-store-microservice/order-service/models"
	"online-store-microservice/order-service/repository"
	orderpb "online-store-microservice/proto/order"
)

var (
	ErrInvalidUserID    = errors.New("invalid user_id")
	ErrInvalidProduct   = errors.New("product_name is required")
	ErrInvalidQuantity  = errors.New("quantity must be greater than 0")
	ErrInvalidPrice     = errors.New("total_price must be greater than 0")
	ErrInvalidOrderID   = errors.New("invalid order id")
	ErrInvalidUserParam = errors.New("invalid user id")
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error)
	GetOrderByID(ctx context.Context, req *orderpb.GetOrderByIdRequest) (*orderpb.GetOrderByIdResponse, error)
	GetOrdersByUserID(ctx context.Context, req *orderpb.GetOrdersByUserIdRequest) (*orderpb.GetOrdersByUserIdResponse, error)
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	if _, err := uuid.Parse(req.UserId); err != nil {
		return nil, ErrInvalidUserID
	}
	if strings.TrimSpace(req.ProductName) == "" {
		return nil, ErrInvalidProduct
	}
	if req.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}
	if req.TotalPrice <= 0 {
		return nil, ErrInvalidPrice
	}

	now := time.Now().UTC()
	order := &models.Order{
		ID:          uuid.NewString(),
		UserID:      req.UserId,
		ProductName: strings.TrimSpace(req.ProductName),
		Quantity:    req.Quantity,
		TotalPrice:  req.TotalPrice,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{Order: toPBOrder(order)}, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, req *orderpb.GetOrderByIdRequest) (*orderpb.GetOrderByIdResponse, error) {
	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, ErrInvalidOrderID
	}

	order, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &orderpb.GetOrderByIdResponse{Order: toPBOrder(order)}, nil
}

func (s *orderService) GetOrdersByUserID(ctx context.Context, req *orderpb.GetOrdersByUserIdRequest) (*orderpb.GetOrdersByUserIdResponse, error) {
	if _, err := uuid.Parse(req.UserId); err != nil {
		return nil, ErrInvalidUserParam
	}

	orders, err := s.repo.GetByUserID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	resp := &orderpb.GetOrdersByUserIdResponse{Orders: make([]*orderpb.OrderData, 0, len(orders))}
	for i := range orders {
		resp.Orders = append(resp.Orders, toPBOrder(&orders[i]))
	}

	return resp, nil
}

func toPBOrder(order *models.Order) *orderpb.OrderData {
	return &orderpb.OrderData{
		Id:          order.ID,
		UserId:      order.UserID,
		ProductName: order.ProductName,
		Quantity:    order.Quantity,
		TotalPrice:  order.TotalPrice,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}
}
