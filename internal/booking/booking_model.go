package booking

import (
	"time"
)

// Reservation represents a reservation entity
type Reservation struct {
	ID                 int                    `json:"id" db:"id"`
	OrganizationID     int                    `json:"organization_id" db:"organization_id"`
	PropertyID         int                    `json:"property_id" db:"property_id"`
	CustomerID         int                    `json:"customer_id" db:"customer_id"`
	CheckInDate        time.Time              `json:"check_in_date" db:"check_in_date"`
	CheckOutDate       time.Time              `json:"check_out_date" db:"check_out_date"`
	Status             string                 `json:"status" db:"status"`
	TotalPrice         float64                `json:"total_price" db:"total_price"`
	PaymentURL         string                 `json:"payment_url" db:"payment_url"`
	PriceElements      map[string]interface{} `json:"price_elements" db:"price_elements"`
	NoOfGuests         int                    `json:"no_of_guests" db:"no_of_guests"`
	GuestData          map[string]interface{} `json:"guest_data" db:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests" db:"additional_requests"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"update_at"`
}

// ReservationRequest represents the reservation creation/update request
type ReservationRequest struct {
	OrganizationID     int                    `json:"organization_id" binding:"required" example:"1"`
	PropertyID         int                    `json:"property_id" binding:"required" example:"10"`
	CustomerID         int                    `json:"customer_id" binding:"required" example:"100"`
	CheckInDate        time.Time              `json:"check_in_date" binding:"required" example:"2024-12-20T15:00:00Z"`
	CheckOutDate       time.Time              `json:"check_out_date" binding:"required" example:"2024-12-25T11:00:00Z"`
	NoOfGuests         int                    `json:"no_of_guests" binding:"required,min=1" example:"2"`
	TotalPrice         float64                `json:"total_price" binding:"required,min=0" example:"500.00"`
	PriceElements      map[string]interface{} `json:"price_elements"`
	GuestData          map[string]interface{} `json:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests"`
	Status             string                 `json:"status" example:"CREATED"`
}

// ReservationResponse represents the detailed reservation response
type ReservationResponse struct {
	ID                 int                    `json:"id" example:"1"`
	OrganizationID     int                    `json:"organization_id" example:"1"`
	PropertyID         int                    `json:"property_id" example:"10"`
	CustomerID         int                    `json:"customer_id" example:"100"`
	CheckInDate        time.Time              `json:"check_in_date" example:"2024-12-20T15:00:00Z"`
	CheckOutDate       time.Time              `json:"check_out_date" example:"2024-12-25T11:00:00Z"`
	Status             string                 `json:"status" example:"CREATED" enums:"CREATED CONFIRMED PAYMENT_REQUIRED REJECTED CANCELLED COMPLETED"`
	TotalPrice         float64                `json:"total_price" example:"500.00"`
	PriceElements      map[string]interface{} `json:"price_elements"`
	NoOfGuests         int                    `json:"no_of_guests" example:"2"`
	GuestData          map[string]interface{} `json:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests"`
	CreatedAt          time.Time              `json:"created_at" example:"2024-12-01T09:00:00Z"`
	UpdatedAt          time.Time              `json:"updated_at" example:"2024-12-01T09:00:00Z"`
}

// ErrorResponse (ostane nespremenjen)
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"The provided data is invalid"`
}

// StatusUpdateRequest (ostane nespremenjen)
type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=CREATED CONFIRMED PAYMENT_REQUIRED REJECTED CANCELLED COMPLETED" example:"CONFIRMED"`
}

// ToResponse (ostane nespremenjen)
func (r *Reservation) ToResponse() *ReservationResponse {
	return &ReservationResponse{
		ID:                 r.ID,
		OrganizationID:     r.OrganizationID,
		PropertyID:         r.PropertyID,
		CustomerID:         r.CustomerID,
		CheckInDate:        r.CheckInDate,
		CheckOutDate:       r.CheckOutDate,
		Status:             r.Status,
		TotalPrice:         r.TotalPrice,
		PriceElements:      r.PriceElements,
		NoOfGuests:         r.NoOfGuests,
		GuestData:          r.GuestData,
		AdditionalRequests: r.AdditionalRequests,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}
