package booking

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"time"
)

// ReservationService handles business logic for reservations
type ReservationService struct {
	repo *ReservationRepository
}

// GetReservationService creates a new ReservationService
func GetReservationService(repo *ReservationRepository) *ReservationService {
	return &ReservationService{
		repo: repo,
	}
}

// GetReservations returns all reservations
func (s *ReservationService) GetReservations() ([]Reservation, error) {
	reservations, err := s.repo.GetReservations()
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationByID returns a reservation by ID
func (s *ReservationService) GetReservationByID(id int) (*Reservation, error) {
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}

	if reservation == nil {
		return nil, errors.New("reservation not found")
	}

	return reservation, nil
}

// CreateReservation creates a new reservation
func (s *ReservationService) CreateReservation(req *ReservationRequest) (*Reservation, error) {
	/*	// Validate request
		if err := s.validateReservationRequest(req); err != nil {
			return nil, err
		}*/

	// Check property availability
	hasConflict, err := s.repo.CheckPropertyAvailability(req.PropertyID, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, errors.New("property is not available for the selected dates")
	}

	// Create reservation entity
	reservation := &Reservation{
		ID:                 rand.IntN(1000000),
		OrganizationID:     req.OrganizationID,
		PropertyID:         req.PropertyID,
		CustomerID:         req.CustomerID,
		CheckInDate:        req.CheckInDate,
		CheckOutDate:       req.CheckOutDate,
		Status:             "CREATED",
		TotalPrice:         req.TotalPrice,
		PriceElements:      req.PriceElements,
		NoOfGuests:         req.NoOfGuests,
		GuestData:          req.GuestData,
		AdditionalRequests: req.AdditionalRequests,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Initialize empty maps if nil
	if reservation.PriceElements == nil {
		reservation.PriceElements = make(map[string]interface{})
	}
	if reservation.GuestData == nil {
		reservation.GuestData = make(map[string]interface{})
	}
	if reservation.AdditionalRequests == nil {
		reservation.AdditionalRequests = make(map[string]interface{})
	}

	// Save to repository
	createdReservation, err := s.repo.CreateReservation(reservation)
	if err != nil {
		return nil, err
	}

	// 5. Call Payment Service
	paymentUrl, err := s.initiatePayment(createdReservation)
	if err != nil {
		// Log error, but you might want to handle failures (e.g., mark as FAILED)
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	createdReservation.PaymentURL = paymentUrl
	createdReservation.Status = "PAYMENT_REQUIRED"
	createdReservation.UpdatedAt = time.Now()

	updatedReservation, err := s.repo.UpdateReservation(createdReservation)
	if err != nil {
		return nil, err
	}

	return updatedReservation, nil
}

func (s *ReservationService) ConfirmPayment(reservationID int) error {
	existing, err := s.repo.GetReservationByID(reservationID)
	if err != nil {
		return fmt.Errorf("could not find reservation %d to confirm: %w", reservationID, err)
	}

	existing.Status = "CONFIRMED"
	existing.UpdatedAt = time.Now()

	_, err = s.repo.UpdateReservation(existing)
	if err != nil {
		return fmt.Errorf("failed to update status for reservation %d: %w", reservationID, err)
	}

	fmt.Printf("Reservation %d status set to 'CONFIMED'\n", reservationID)
	return nil
}

// UpdateReservation updates an existing reservation
func (s *ReservationService) UpdateReservation(id int, req *ReservationRequest) (*Reservation, error) {
	// Validate request
	/*	if err := s.validateReservationRequest(req); err != nil {
		return nil, err
	}*/

	// Check if reservation exists
	existingReservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}
	if existingReservation == nil {
		return nil, errors.New("reservation not found")
	}

	// Check if reservation can be updated
	if existingReservation.Status == "completed" || existingReservation.Status == "cancelled" {
		return nil, errors.New("cannot update a completed or cancelled reservation")
	}

	// Check property availability (excluding current reservation)
	hasConflict, err := s.repo.CheckPropertyAvailabilityExcluding(id, req.PropertyID, req.CheckInDate, req.CheckOutDate)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, errors.New("property is not available for the selected dates")
	}

	// Update reservation fields
	existingReservation.OrganizationID = req.OrganizationID
	existingReservation.PropertyID = req.PropertyID
	existingReservation.CustomerID = req.CustomerID
	existingReservation.CheckInDate = req.CheckInDate
	existingReservation.CheckOutDate = req.CheckOutDate
	existingReservation.TotalPrice = req.TotalPrice
	existingReservation.PriceElements = req.PriceElements
	existingReservation.NoOfGuests = req.NoOfGuests
	existingReservation.GuestData = req.GuestData
	existingReservation.AdditionalRequests = req.AdditionalRequests
	existingReservation.UpdatedAt = time.Now()

	// Initialize empty maps if nil
	if existingReservation.PriceElements == nil {
		existingReservation.PriceElements = make(map[string]interface{})
	}
	if existingReservation.GuestData == nil {
		existingReservation.GuestData = make(map[string]interface{})
	}
	if existingReservation.AdditionalRequests == nil {
		existingReservation.AdditionalRequests = make(map[string]interface{})
	}

	// Save updates
	updatedReservation, err := s.repo.UpdateReservation(existingReservation)
	if err != nil {
		return nil, err
	}

	return updatedReservation, nil
}

func (s *ReservationService) initiatePayment(res *Reservation) (string, error) {
	paymentReq := map[string]interface{}{
		"organizationId": res.OrganizationID,
		"reservationId":  res.ID,
		"customerId":     res.CustomerID,
		"amount":         res.TotalPrice,
	}

	body, err := json.Marshal(paymentReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("https://hostflow.software/payment/payments", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("payment service returned status: %d", resp.StatusCode)
	}

	// Read the string response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

// DeleteReservation deletes a reservation by ID
func (s *ReservationService) DeleteReservation(id int) error {
	// Check if reservation exists
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	// Business rule: cannot delete completed reservations
	if reservation.Status == "completed" {
		return errors.New("cannot delete a completed reservation")
	}

	// Delete the reservation
	err = s.repo.DeleteReservation(id)
	if err != nil {
		return err
	}

	return nil
}

// UpdateReservationStatus updates only the status of a reservation
func (s *ReservationService) UpdateReservationStatus(id int, status string) (*Reservation, error) {
	// Validate status
	validStatuses := map[string]bool{
		"pending": true, "confirmed": true, "checked_in": true,
		"checked_out": true, "cancelled": true, "completed": true, "rejected": true,
	}
	if !validStatuses[status] {
		return nil, errors.New("invalid status value")
	}

	// Get existing reservation
	reservation, err := s.repo.GetReservationByID(id)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, errors.New("reservation not found")
	}

	// Validate status transition
	if err := s.validateStatusTransition(reservation.Status, status); err != nil {
		return nil, err
	}

	// Update status
	reservation.Status = status
	reservation.UpdatedAt = time.Now()

	updatedReservation, err := s.repo.UpdateReservation(reservation)
	if err != nil {
		return nil, err
	}

	return updatedReservation, nil
}

// CancelReservation cancels a reservation
func (s *ReservationService) CancelReservation(id int) (*Reservation, error) {
	return s.UpdateReservationStatus(id, "cancelled")
}

// ConfirmReservation confirms a pending reservation
func (s *ReservationService) ConfirmReservation(id int) (*Reservation, error) {
	return s.UpdateReservationStatus(id, "confirmed")
}

// CheckInReservation marks a reservation as checked in
func (s *ReservationService) CheckInReservation(id int) (*Reservation, error) {
	return s.UpdateReservationStatus(id, "checked_in")
}

// CheckOutReservation marks a reservation as checked out
func (s *ReservationService) CheckOutReservation(id int) (*Reservation, error) {
	return s.UpdateReservationStatus(id, "checked_out")
}

// GetReservationsByCustomer returns all reservations for a customer
func (s *ReservationService) GetReservationsByCustomer(customerID int) ([]Reservation, error) {
	reservations, err := s.repo.GetReservationsByCustomer(customerID)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByProperty returns all reservations for a property
func (s *ReservationService) GetReservationsByProperty(propertyID int) ([]Reservation, error) {
	reservations, err := s.repo.GetReservationsByProperty(propertyID)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByOrganization returns all reservations for an organization
func (s *ReservationService) GetReservationsByOrganization(organizationID int) ([]Reservation, error) {
	reservations, err := s.repo.GetReservationsByOrganization(organizationID)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

/*// validateReservationRequest validates the reservation request
func (s *ReservationService) validateReservationRequest(req *ReservationRequest) error {
	if req.OrganizationID == uuid.Nil {
		return errors.New("organization ID is required")
	}
	if req.PropertyID == uuid.Nil {
		return errors.New("property ID is required")
	}
	if req.CheckInDate.IsZero() {
		return errors.New("check-in date is required")
	}
	if req.CheckOutDate.IsZero() {
		return errors.New("check-out date is required")
	}
	if req.CheckOutDate.Before(req.CheckInDate) || req.CheckOutDate.Equal(req.CheckInDate) {
		return errors.New("check-out date must be after check-in date")
	}
	if req.CheckInDate.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("check-in date cannot be in the past")
	}
	if req.NoOfGuests < 1 {
		return errors.New("number of guests must be at least 1")
	}
	if req.TotalPrice < 0 {
		return errors.New("total price cannot be negative")
	}

	// Validate minimum stay (1 night)
	duration := req.CheckOutDate.Sub(req.CheckInDate)
	if duration < 24*time.Hour {
		return errors.New("minimum stay is 1 night")
	}

	return nil
}*/

// validateStatusTransition validates if a status transition is allowed
func (s *ReservationService) validateStatusTransition(currentStatus, newStatus string) error {
	// Define allowed transitions
	allowedTransitions := map[string][]string{
		"pending":     {"confirmed", "cancelled", "rejected"},
		"confirmed":   {"checked_in", "cancelled"},
		"checked_in":  {"checked_out", "cancelled"},
		"checked_out": {"completed"},
		"cancelled":   {}, // Cannot transition from cancelled
		"completed":   {}, // Cannot transition from completed
		"rejected":    {}, // Cannot transition from rejected
	}

	if currentStatus == newStatus {
		return nil // Same status is allowed
	}

	allowed, exists := allowedTransitions[currentStatus]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, allowedStatus := range allowed {
		if allowedStatus == newStatus {
			return nil
		}
	}

	return errors.New("status transition not allowed from " + currentStatus + " to " + newStatus)
}
