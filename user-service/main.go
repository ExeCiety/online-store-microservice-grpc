package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"online-store-microservice/pkg/grpcjson"
	pkglog "online-store-microservice/pkg/logger"
	userpb "online-store-microservice/proto/user"
	"online-store-microservice/user-service/config"
	"online-store-microservice/user-service/repository"
	"online-store-microservice/user-service/server"
	"online-store-microservice/user-service/service"
)

func main() {
	grpcjson.Register()
	log := pkglog.New("[user-service]")
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	grpcSrv := server.NewGRPCServer(svc, log)

	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(loggingInterceptor(log)))
	userpb.RegisterUserServiceServer(s, grpcSrv)

	go func() {
		log.Printf("gRPC server listening on :%s", cfg.GRPCPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("serve: %v", err)
		}
	}()

	shutdown(log, s)
}

func loggingInterceptor(log interface{ Printf(string, ...interface{}) }) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		log.Printf("method=%s duration=%s err=%v", info.FullMethod, time.Since(start), err)
		return resp, err
	}
}

func shutdown(log interface{ Printf(string, ...interface{}) }, s *grpc.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Printf("shutting down user service")
	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Printf("user service stopped")
	case <-time.After(10 * time.Second):
		log.Printf("force stop user service after timeout")
		s.Stop()
	}
	fmt.Println()
}
