package grpc_clients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"online-store-microservice/pkg/grpcjson"
	userpb "online-store-microservice/proto/user"
)

type UserClient struct {
	conn   *grpc.ClientConn
	Client userpb.UserServiceClient
}

func NewUserClient(addr string) (*UserClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcjson.Codec{})),
	)
	if err != nil {
		return nil, err
	}

	return &UserClient{conn: conn, Client: userpb.NewUserServiceClient(conn)}, nil
}

func (c *UserClient) Close() error {
	return c.conn.Close()
}

func (c *UserClient) TimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
