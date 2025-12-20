package booking

import (
	"time"

	"github.com/google/uuid"
)

// Reservation represents a reservation entity
type Reservation struct {
	ID                 uuid.UUID              `json:"id" db:"id"`
	OrganizationID     uuid.UUID              `json:"organization_id" db:"organization_id"`
	PropertyID         uuid.UUID              `json:"property_id" db:"property_id"`
	CustomerID         uuid.UUID              `json:"customer_id" db:"customer_id"`
	CheckInDate        time.Time              `json:"check_in_date" db:"check_in_date"`
	Status             string                 `json:"status" db:"status"`
	TotalPrice         float64                `json:"total_price" db:"total_price"`
	PriceElements      map[string]interface{} `json:"price_elements" db:"price_elements"`
	NoOfGuests         int                    `json:"no_of_guests" db:"no_of_guests"`
	GuestData          map[string]interface{} `json:"guest_data" db:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests" db:"additional_requests"`
	CheckOutDate       time.Time              `json:"check_out_date" db:"check_out_date"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"update_at"`
}

// ReservationRequest represents the reservation creation/update request
type ReservationRequest struct {
	OrganizationID     uuid.UUID              `json:"organization_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	PropertyID         uuid.UUID              `json:"property_id" binding:"required" example:"660e8400-e29b-41d4-a716-446655440000"`
	CustomerID         uuid.UUID              `json:"customer_id" binding:"required" example:"770e8400-e29b-41d4-a716-446655440000"`
	CheckInDate        time.Time              `json:"check_in_date" binding:"required" example:"2024-12-20T15:00:00Z"`
	CheckOutDate       time.Time              `json:"check_out_date" binding:"required" example:"2024-12-25T11:00:00Z"`
	NoOfGuests         int                    `json:"no_of_guests" binding:"required,min=1" example:"2"`
	TotalPrice         float64                `json:"total_price" binding:"required,min=0" example:"500.00"`
	PriceElements      map[string]interface{} `json:"price_elements"`
	GuestData          map[string]interface{} `json:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests"`
}

// ReservationResponse represents the detailed reservation response
type ReservationResponse struct {
	ID                 uuid.UUID              `json:"id" example:"880e8400-e29b-41d4-a716-446655440000"`
	OrganizationID     uuid.UUID              `json:"organization_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	PropertyID         uuid.UUID              `json:"property_id" example:"660e8400-e29b-41d4-a716-446655440000"`
	CustomerID         uuid.UUID              `json:"customer_id" example:"770e8400-e29b-41d4-a716-446655440000"`
	CheckInDate        time.Time              `json:"check_in_date" example:"2024-12-20T15:00:00Z"`
	CheckOutDate       time.Time              `json:"check_out_date" example:"2024-12-25T11:00:00Z"`
	Status             string                 `json:"status" example:"confirmed" enums:"pending,confirmed,checked_in,checked_out,cancelled,completed,rejected"`
	TotalPrice         float64                `json:"total_price" example:"500.00"`
	PriceElements      map[string]interface{} `json:"price_elements"`
	NoOfGuests         int                    `json:"no_of_guests" example:"2"`
	GuestData          map[string]interface{} `json:"guest_data"`
	AdditionalRequests map[string]interface{} `json:"additional_requests"`
	CreatedAt          time.Time              `json:"created_at" example:"2024-12-01T09:00:00Z"`
	UpdatedAt          time.Time              `json:"updated_at" example:"2024-12-01T09:00:00Z"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"The provided data is invalid"`
}

// StatusUpdateRequest represents a status update request
type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed checked_in checked_out cancelled completed rejected" example:"confirmed"`
}

// Helper function to convert Reservation to ReservationResponse
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
