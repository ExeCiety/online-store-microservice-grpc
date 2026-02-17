package server

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"online-store-microservice/order-service/service"
	orderpb "online-store-microservice/proto/order"
)

type GRPCServer struct {
	orderpb.UnimplementedOrderServiceServer
	service service.OrderService
}

func NewGRPCServer(svc service.OrderService) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	resp, err := s.service.CreateOrder(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (s *GRPCServer) GetOrderById(ctx context.Context, req *orderpb.GetOrderByIdRequest) (*orderpb.GetOrderByIdResponse, error) {
	resp, err := s.service.GetOrderByID(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func (s *GRPCServer) GetOrdersByUserId(ctx context.Context, req *orderpb.GetOrdersByUserIdRequest) (*orderpb.GetOrdersByUserIdResponse, error) {
	resp, err := s.service.GetOrdersByUserID(ctx, req)
	if err != nil {
		return nil, mapError(err)
	}
	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidUserID),
		errors.Is(err, service.ErrInvalidProduct),
		errors.Is(err, service.ErrInvalidQuantity),
		errors.Is(err, service.ErrInvalidPrice),
		errors.Is(err, service.ErrInvalidOrderID),
		errors.Is(err, service.ErrInvalidUserParam):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		return status.Error(codes.NotFound, "order not found")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
