package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"online-store-microservice/api-gateway/grpc_clients"
	"online-store-microservice/pkg/response"
	userpb "online-store-microservice/proto/user"
)

type UserHandler struct {
	client *grpc_clients.UserClient
}

func NewUserHandler(client *grpc_clients.UserClient) *UserHandler {
	return &UserHandler{client: client}
}

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required,min=2"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.Register(ctx, &userpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to register user", msg)
		return
	}

	response.OK(c, http.StatusCreated, "user registered", resp.User)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.Login(ctx, &userpb.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to login", msg)
		return
	}

	response.OK(c, http.StatusOK, "login successful", gin.H{"user": resp.User, "token": resp.Token})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Fail(c, http.StatusBadRequest, "id is required", nil)
		return
	}

	ctx, cancel := h.client.TimeoutContext()
	defer cancel()

	resp, err := h.client.Client.GetUserById(ctx, &userpb.GetUserByIdRequest{Id: id})
	if err != nil {
		code, msg := grpcToHTTP(err)
		response.Fail(c, code, "failed to get user", msg)
		return
	}

	response.OK(c, http.StatusOK, "user fetched", resp.User)
}
