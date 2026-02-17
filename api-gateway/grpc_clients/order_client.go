package grpc_clients

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"online-store-microservice/pkg/grpcjson"
	orderpb "online-store-microservice/proto/order"
)

type OrderClient struct {
	conn   *grpc.ClientConn
	Client orderpb.OrderServiceClient
}

func NewOrderClient(addr string) (*OrderClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.ForceCodec(grpcjson.Codec{})),
	)
	if err != nil {
		return nil, err
	}

	return &OrderClient{conn: conn, Client: orderpb.NewOrderServiceClient(conn)}, nil
}

func (c *OrderClient) Close() error {
	return c.conn.Close()
}

func (c *OrderClient) TimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
