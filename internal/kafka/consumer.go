package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.uber.org/fx"
)

func NewKafkaReader() (*kafka.Reader, error) {
	// 1. Setup SCRAM-SHA-256 Mechanism
	mechanism, err := scram.Mechanism(scram.SHA256,
		os.Getenv("KAFKA_USER"),
		os.Getenv("KAFKA_PASSWORD"),
	)
	if err != nil {
		return nil, err
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   "booking.payments",
		GroupID: "communication-service-group",
		Dialer:  dialer,
	}), nil
}

func RegisterKafkaHooks(lifecycle fx.Lifecycle, reader *kafka.Reader) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Kafka Consumer starting on topic: booking.payments")

			go func() {
				for {
					m, err := reader.ReadMessage(context.Background())
					if err != nil {
						fmt.Printf("Kafka error: %v\n", err)
						return
					}

					var envelope MessageEnvelope
					if err := json.Unmarshal(m.Value, &envelope); err != nil {
						fmt.Printf("Failed to parse message: %v\n", err)
						continue
					}

					// Logic check: Only process if it's a PaymentAction and succeeded
					if envelope.MessageType == "PaymentAction" && envelope.Payload.StripeStatus == "succeeded" {
						fmt.Printf("Payment Succeeded for Reservation: %d\n", envelope.Payload.ReservationId)

						// TODO dodaj klic v booking service
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Kafka Consumer shutting down...")
			return reader.Close()
		},
	})
}

var Module = fx.Module("kafka",
	fx.Provide(NewKafkaReader),
	fx.Invoke(RegisterKafkaHooks),
)
