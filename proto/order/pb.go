package orderpb

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderData struct {
	Id          string  `json:"id"`
	UserId      string  `json:"user_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateOrderRequest struct {
	UserId      string  `json:"user_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	TotalPrice  float64 `json:"total_price"`
}

type CreateOrderResponse struct {
	Order *OrderData `json:"order"`
}

type GetOrderByIdRequest struct {
	Id string `json:"id"`
}

type GetOrderByIdResponse struct {
	Order *OrderData `json:"order"`
}

type GetOrdersByUserIdRequest struct {
	UserId string `json:"user_id"`
}

type GetOrdersByUserIdResponse struct {
	Orders []*OrderData `json:"orders"`
}

type OrderServiceClient interface {
	CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error)
	GetOrderById(ctx context.Context, in *GetOrderByIdRequest, opts ...grpc.CallOption) (*GetOrderByIdResponse, error)
	GetOrdersByUserId(ctx context.Context, in *GetOrdersByUserIdRequest, opts ...grpc.CallOption) (*GetOrdersByUserIdResponse, error)
}

type orderServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOrderServiceClient(cc grpc.ClientConnInterface) OrderServiceClient {
	return &orderServiceClient{cc: cc}
}

func (c *orderServiceClient) CreateOrder(ctx context.Context, in *CreateOrderRequest, opts ...grpc.CallOption) (*CreateOrderResponse, error) {
	out := new(CreateOrderResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/CreateOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetOrderById(ctx context.Context, in *GetOrderByIdRequest, opts ...grpc.CallOption) (*GetOrderByIdResponse, error) {
	out := new(GetOrderByIdResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/GetOrderById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *orderServiceClient) GetOrdersByUserId(ctx context.Context, in *GetOrdersByUserIdRequest, opts ...grpc.CallOption) (*GetOrdersByUserIdResponse, error) {
	out := new(GetOrdersByUserIdResponse)
	err := c.cc.Invoke(ctx, "/order.OrderService/GetOrdersByUserId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type OrderServiceServer interface {
	CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error)
	GetOrderById(context.Context, *GetOrderByIdRequest) (*GetOrderByIdResponse, error)
	GetOrdersByUserId(context.Context, *GetOrdersByUserIdRequest) (*GetOrdersByUserIdResponse, error)
	mustEmbedUnimplementedOrderServiceServer()
}

type UnimplementedOrderServiceServer struct{}

func (UnimplementedOrderServiceServer) CreateOrder(context.Context, *CreateOrderRequest) (*CreateOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}

func (UnimplementedOrderServiceServer) GetOrderById(context.Context, *GetOrderByIdRequest) (*GetOrderByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrderById not implemented")
}

func (UnimplementedOrderServiceServer) GetOrdersByUserId(context.Context, *GetOrdersByUserIdRequest) (*GetOrdersByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOrdersByUserId not implemented")
}

func (UnimplementedOrderServiceServer) mustEmbedUnimplementedOrderServiceServer() {}

func RegisterOrderServiceServer(s grpc.ServiceRegistrar, srv OrderServiceServer) {
	s.RegisterService(&OrderService_ServiceDesc, srv)
}

func _OrderService_CreateOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).CreateOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/order.OrderService/CreateOrder"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).CreateOrder(ctx, req.(*CreateOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetOrderById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrderByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetOrderById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/order.OrderService/GetOrderById"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetOrderById(ctx, req.(*GetOrderByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OrderService_GetOrdersByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOrdersByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrderServiceServer).GetOrdersByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/order.OrderService/GetOrdersByUserId"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrderServiceServer).GetOrdersByUserId(ctx, req.(*GetOrdersByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var OrderService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "order.OrderService",
	HandlerType: (*OrderServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "CreateOrder", Handler: _OrderService_CreateOrder_Handler},
		{MethodName: "GetOrderById", Handler: _OrderService_GetOrderById_Handler},
		{MethodName: "GetOrdersByUserId", Handler: _OrderService_GetOrdersByUserId_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "order.proto",
}
