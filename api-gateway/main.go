package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"online-store-microservice/api-gateway/config"
	"online-store-microservice/api-gateway/grpc_clients"
	"online-store-microservice/api-gateway/handlers"
	"online-store-microservice/api-gateway/middleware"
	pkglog "online-store-microservice/pkg/logger"
	"online-store-microservice/pkg/response"
)

func main() {
	log := pkglog.New("[api-gateway]")
	cfg := config.Load()

	userClient, err := grpc_clients.NewUserClient(cfg.UserServiceURL)
	if err != nil {
		log.Fatalf("connect user service: %v", err)
	}
	defer userClient.Close()

	orderClient, err := grpc_clients.NewOrderClient(cfg.OrderServiceURL)
	if err != nil {
		log.Fatalf("connect order service: %v", err)
	}
	defer orderClient.Close()

	userHandler := handlers.NewUserHandler(userClient)
	orderHandler := handlers.NewOrderHandler(orderClient)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		rid := param.Request.Header.Get("X-Request-ID")
		return fmt.Sprintf("[api-gateway] %s | %3d | %13v | %15s | %s %s | req_id=%s\n",
			param.TimeStamp.UTC().Format(time.RFC3339),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			rid,
		)
	}))

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, http.StatusOK, "ok", gin.H{"service": "api-gateway"})
	})

	api := r.Group("/api")
	api.POST("/register", userHandler.Register)
	api.POST("/login", userHandler.Login)
	api.GET("/users/:id", userHandler.GetByID)
	api.POST("/orders", orderHandler.Create)
	api.GET("/orders/:id", orderHandler.GetByID)
	api.GET("/users/:id/orders", orderHandler.GetByUserID)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("http server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	shutdown(log, srv)
}

func shutdown(log interface{ Printf(string, ...interface{}) }, srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Printf("shutting down api gateway")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
