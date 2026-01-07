package kafka

type PaymentAction struct {
	PaymentId             int64   `json:"paymentId"`
	OrganizationId        int64   `json:"organizationId"`
	ReservationId         int64   `json:"reservationId"` // Matches your DB int8
	Amount                float64 `json:"amount"`
	StripePaymentIntentId string  `json:"stripePaymentIntentId"`
	StripeStatus          string  `json:"stripeStatus"`
	PaidAtUtc             string  `json:"paidAtUtc"`
}

type MessageEnvelope struct {
	MessageId     string        `json:"messageId"`
	MessageType   string        `json:"messageType"`
	OccurredAt    string        `json:"occurredAt"`
	Payload       PaymentAction `json:"payload"`
	SchemaVersion int           `json:"schemaVersion"`
}
