package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"online-store-microservice/api-gateway/grpc_clients"
	"online-store-microservice/api-gateway/handlers"
	"online-store-microservice/pkg/response"
	orderpb "online-store-microservice/proto/order"
	userpb "online-store-microservice/proto/user"
)

type fakeUserServiceClient struct {
	registerFn    func(context.Context, *userpb.RegisterRequest, ...grpc.CallOption) (*userpb.RegisterResponse, error)
	loginFn       func(context.Context, *userpb.LoginRequest, ...grpc.CallOption) (*userpb.LoginResponse, error)
	getUserByIDFn func(context.Context, *userpb.GetUserByIdRequest, ...grpc.CallOption) (*userpb.GetUserByIdResponse, error)
}

func (f *fakeUserServiceClient) Register(ctx context.Context, req *userpb.RegisterRequest, opts ...grpc.CallOption) (*userpb.RegisterResponse, error) {
	return f.registerFn(ctx, req, opts...)
}

func (f *fakeUserServiceClient) Login(ctx context.Context, req *userpb.LoginRequest, opts ...grpc.CallOption) (*userpb.LoginResponse, error) {
	return f.loginFn(ctx, req, opts...)
}

func (f *fakeUserServiceClient) GetUserById(ctx context.Context, req *userpb.GetUserByIdRequest, opts ...grpc.CallOption) (*userpb.GetUserByIdResponse, error) {
	return f.getUserByIDFn(ctx, req, opts...)
}

type fakeOrderServiceClient struct {
	createOrderFn       func(context.Context, *orderpb.CreateOrderRequest, ...grpc.CallOption) (*orderpb.CreateOrderResponse, error)
	getOrderByIDFn      func(context.Context, *orderpb.GetOrderByIdRequest, ...grpc.CallOption) (*orderpb.GetOrderByIdResponse, error)
	getOrdersByUserIDFn func(context.Context, *orderpb.GetOrdersByUserIdRequest, ...grpc.CallOption) (*orderpb.GetOrdersByUserIdResponse, error)
}

func (f *fakeOrderServiceClient) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest, opts ...grpc.CallOption) (*orderpb.CreateOrderResponse, error) {
	return f.createOrderFn(ctx, req, opts...)
}

func (f *fakeOrderServiceClient) GetOrderById(ctx context.Context, req *orderpb.GetOrderByIdRequest, opts ...grpc.CallOption) (*orderpb.GetOrderByIdResponse, error) {
	return f.getOrderByIDFn(ctx, req, opts...)
}

func (f *fakeOrderServiceClient) GetOrdersByUserId(ctx context.Context, req *orderpb.GetOrdersByUserIdRequest, opts ...grpc.CallOption) (*orderpb.GetOrdersByUserIdResponse, error) {
	return f.getOrdersByUserIDFn(ctx, req, opts...)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC().Format(time.RFC3339)

	fakeUser := &fakeUserServiceClient{
		registerFn: func(context.Context, *userpb.RegisterRequest, ...grpc.CallOption) (*userpb.RegisterResponse, error) {
			return &userpb.RegisterResponse{User: &userpb.UserData{Id: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", Email: "user@example.com", Name: "John Doe", CreatedAt: now, UpdatedAt: now}}, nil
		},
		loginFn: func(context.Context, *userpb.LoginRequest, ...grpc.CallOption) (*userpb.LoginResponse, error) {
			return &userpb.LoginResponse{User: &userpb.UserData{Id: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", Email: "user@example.com", Name: "John Doe", CreatedAt: now, UpdatedAt: now}, Token: "dummy-token"}, nil
		},
		getUserByIDFn: func(context.Context, *userpb.GetUserByIdRequest, ...grpc.CallOption) (*userpb.GetUserByIdResponse, error) {
			return &userpb.GetUserByIdResponse{User: &userpb.UserData{Id: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", Email: "user@example.com", Name: "John Doe", CreatedAt: now, UpdatedAt: now}}, nil
		},
	}

	fakeOrder := &fakeOrderServiceClient{
		createOrderFn: func(context.Context, *orderpb.CreateOrderRequest, ...grpc.CallOption) (*orderpb.CreateOrderResponse, error) {
			return &orderpb.CreateOrderResponse{Order: &orderpb.OrderData{Id: "8f328abb-4ae4-493b-a460-a63f1206b2f3", UserId: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", ProductName: "Laptop", Quantity: 1, TotalPrice: 15000000, Status: "pending", CreatedAt: now, UpdatedAt: now}}, nil
		},
		getOrderByIDFn: func(context.Context, *orderpb.GetOrderByIdRequest, ...grpc.CallOption) (*orderpb.GetOrderByIdResponse, error) {
			return &orderpb.GetOrderByIdResponse{Order: &orderpb.OrderData{Id: "8f328abb-4ae4-493b-a460-a63f1206b2f3", UserId: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", ProductName: "Laptop", Quantity: 1, TotalPrice: 15000000, Status: "pending", CreatedAt: now, UpdatedAt: now}}, nil
		},
		getOrdersByUserIDFn: func(context.Context, *orderpb.GetOrdersByUserIdRequest, ...grpc.CallOption) (*orderpb.GetOrdersByUserIdResponse, error) {
			return &orderpb.GetOrdersByUserIdResponse{Orders: []*orderpb.OrderData{{Id: "8f328abb-4ae4-493b-a460-a63f1206b2f3", UserId: "4e427d78-58c5-4f78-bfc1-e2c196e0b506", ProductName: "Laptop", Quantity: 1, TotalPrice: 15000000, Status: "pending", CreatedAt: now, UpdatedAt: now}}}, nil
		},
	}

	userHandler := handlers.NewUserHandler(&grpc_clients.UserClient{Client: fakeUser})
	orderHandler := handlers.NewOrderHandler(&grpc_clients.OrderClient{Client: fakeOrder})

	r := gin.New()
	r.GET("/health", func(c *gin.Context) { response.OK(c, http.StatusOK, "ok", gin.H{"service": "api-gateway"}) })

	api := r.Group("/api")
	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)
	api.GET("/users/:id", userHandler.GetByID)
	api.POST("/orders", orderHandler.Create)
	api.GET("/orders/:id", orderHandler.GetByID)
	api.GET("/users/:id/orders", orderHandler.GetByUserID)

	return r
}

func doRequest(r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	var payload []byte
	if body != nil {
		payload, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewBuffer(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
