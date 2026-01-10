package communication

import (
	"context"
	"log"

	pb "hostflow/booking-service/internal/communication/proto"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Client pb.CommunicationServiceClient

func NewCustomerClient(lc fx.Lifecycle) (pb.CommunicationServiceClient, error) {
	conn, err := grpc.Dial(
		"communication-service-dev:50051",
		grpc.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})

	client := pb.NewCommunicationServiceClient(conn)
	log.Println("connected to Customer gRPC service")
	return client, nil
}

var Module = fx.Provide(NewCustomerClient, NewController)
