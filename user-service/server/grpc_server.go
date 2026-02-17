package server

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	userpb "online-store-microservice/proto/user"
	"online-store-microservice/user-service/service"
)

type GRPCServer struct {
	userpb.UnimplementedUserServiceServer
	service service.UserService
	logger  *log.Logger
}

func NewGRPCServer(svc service.UserService, logger *log.Logger) *GRPCServer {
	return &GRPCServer{service: svc, logger: logger}
}

func (s *GRPCServer) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	resp, err := s.service.Register(ctx, req)
	if err != nil {
		s.logger.Printf("register failed: %v", err)
		return nil, mapError(err)
	}
	return resp, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	resp, err := s.service.Login(ctx, req)
	if err != nil {
		s.logger.Printf("login failed: %v", err)
		return nil, mapError(err)
	}
	return resp, nil
}

func (s *GRPCServer) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error) {
	resp, err := s.service.GetUserByID(ctx, req)
	if err != nil {
		s.logger.Printf("get user by id failed: %v", err)
		return nil, mapError(err)
	}
	return resp, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidName),
		errors.Is(err, service.ErrInvalidPassword):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, service.ErrEmailAlreadyUsed):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, service.ErrInvalidCredential):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, gorm.ErrRecordNotFound):
		return status.Error(codes.NotFound, "user not found")
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
