package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"online-store-microservice/api-gateway/grpc_clients"
	"online-store-microservice/pkg/response"
	orderpb "online-store-microservice/proto/order"
)

type OrderHandler struct {
	client *grpc_clients.OrderClient
}

func NewOrderHandler(client *grpc_clients.OrderClient) *OrderHandler {
	return &OrderHandler{client: client}
}

type createOrderRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	Quantity    int32   `json:"quantity" binding:"required,gt=0"`
	TotalPrice  float64 `json:"total_price" binding:"required,gt=0"`
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId:      req.UserID,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		TotalPrice:  req.TotalPrice,
	})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to create order", msg)
		return
	}

	response.OK(c, http.StatusCreated, "order created", resp.Order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "id is required", nil)
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.GetOrderById(ctx, &orderpb.GetOrderByIdRequest{Id: id})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to get order", msg)
		return
	}

	response.OK(c, http.StatusOK, "order fetched", resp.Order)
}

func (h *OrderHandler) GetByUserID(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.Param("id")
	}
	if userID == "" {
		response.Fail(c, http.StatusBadRequest, "userId is required", nil)
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.GetOrdersByUserId(ctx, &orderpb.GetOrdersByUserIdRequest{UserId: userID})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to get user orders", msg)
		return
	}

	response.OK(c, http.StatusOK, "orders fetched", resp.Orders)
}
