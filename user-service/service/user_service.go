package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	userpb "online-store-microservice/proto/user"
	"online-store-microservice/user-service/models"
	"online-store-microservice/user-service/repository"
)

var (
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidName       = errors.New("invalid name")
	ErrInvalidPassword   = errors.New("password must be at least 6 characters")
	ErrEmailAlreadyUsed  = errors.New("email already registered")
	ErrInvalidCredential = errors.New("invalid credentials")
)

type UserService interface {
	Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error)
	Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error)
	GetUserByID(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, ErrInvalidEmail
	}
	if len(req.Name) < 2 {
		return nil, ErrInvalidName
	}
	if len(req.Password) < 6 {
		return nil, ErrInvalidPassword
	}

	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrEmailAlreadyUsed
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &models.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		PasswordHash: string(hashed),
		Name:         req.Name,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &userpb.RegisterResponse{User: toPBUser(user)}, nil
}

func (s *userService) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return nil, ErrInvalidCredential
	}

	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredential
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}

	token := fmt.Sprintf("dummy-token-%s", user.ID)
	return &userpb.LoginResponse{User: toPBUser(user), Token: token}, nil
}

func (s *userService) GetUserByID(ctx context.Context, req *userpb.GetUserByIdRequest) (*userpb.GetUserByIdResponse, error) {
	if strings.TrimSpace(req.Id) == "" {
		return nil, errors.New("id is required")
	}

	user, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userpb.GetUserByIdResponse{User: toPBUser(user)}, nil
}

func toPBUser(user *models.User) *userpb.UserData {
	return &userpb.UserData{
		Id:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}
