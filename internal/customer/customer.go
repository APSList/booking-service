package customer

import (
	"context"
	"log"

	pb "hostflow/booking-service/internal/customer/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Client pb.CustomerServiceClient

func NewCustomerClient(lc fx.Lifecycle) (pb.CustomerServiceClient, error) {
	conn, err := grpc.Dial(
		"customer-service-dev:50051",
		grpc.WithInsecure(), // if no TLS inside cluster
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})

	client := pb.NewCustomerServiceClient(conn)
	log.Println("connected to Customer gRPC service")
	return client, nil
}

var Module = fx.Provide(NewCustomerClient, NewController)
