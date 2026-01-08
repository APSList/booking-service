package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"hostflow/booking-service/internal/booking"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"go.uber.org/fx"
)

func NewKafkaReader() *kafka.Reader {
	mechanism := plain.Mechanism{
		Username: os.Getenv("KAFKA_USER"),
		Password: os.Getenv("KAFKA_PASSWORD"),
	}

	//sharedTransport := &kafka.Transport{
	//	SASL: mechanism,
	//	TLS:  &tls.Config{}, // This enables the "SSL" part of SASL_SSL
	//}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   os.Getenv("KAFKA_TOPIC"),
		GroupID: "communication-service-group",
		Dialer: &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		},
		// Replicating Confluent's session timeout behavior
		ReadBatchTimeout: 10 * time.Second,
	})
}

func RegisterKafkaHooks(lifecycle fx.Lifecycle, reader *kafka.Reader, reservationService *booking.ReservationService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Kafka Consumer starting...")
			go func() {
				for {
					m, err := reader.ReadMessage(context.Background())
					if err != nil {
						fmt.Printf("Kafka error: %v\n", err)
						return
					}

					fmt.Printf("Received message: %s\n", string(m.Value))
					var envelope MessageEnvelope
					if err := json.Unmarshal(m.Value, &envelope); err != nil {
						fmt.Printf("Failed to parse message: %v\n", err)
						continue
					}

					// Only process if it's a PaymentAction and succeeded
					if envelope.MessageType == "PaymentAction" && envelope.Payload.StripeStatus == "succeeded" {
						fmt.Printf("Processing successful payment for Reservation: %d\n", envelope.Payload.ReservationId)

						err := reservationService.ConfirmPayment(int(envelope.Payload.ReservationId))
						if err != nil {
							fmt.Printf("Failed to confirm payment for Reservation: %d; err: %s\n", envelope.Payload.ReservationId, err)
						}
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return reader.Close()
		},
	})
}

var Module = fx.Module("kafka",
	fx.Provide(NewKafkaReader),
	fx.Invoke(RegisterKafkaHooks),
)
