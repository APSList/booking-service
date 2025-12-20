package booking

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReservationController handles HTTP requests for reservations
type ReservationController struct {
	service *ReservationService
}

// GetReservationController creates a new controller
func GetReservationController(service *ReservationService) *ReservationController {
	return &ReservationController{
		service: service,
	}
}

// GetReservationsHandler godoc
// @Summary Get all reservations
// @Description Returns a list of all reservations
// @Tags reservations
// @Accept json
// @Produce json
// @Success 200 {array} ReservationResponse
// @Failure 500 {object} ErrorResponse
// @Router /reservations [get]
func (c *ReservationController) GetReservationsHandler(ctx *gin.Context) {
	reservations, err := c.service.GetReservations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch reservations",
			Message: err.Error(),
		})
		return
	}

	// Convert to response format
	response := make([]ReservationResponse, len(reservations))
	for i, r := range reservations {
		response[i] = *r.ToResponse()
	}

	ctx.JSON(http.StatusOK, response)
}

// GetReservationByIDHandler godoc
// @Summary Get reservation by ID
// @Description Get reservation details by reservation ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reservations/{id} [get]
func (c *ReservationController) GetReservationByIDHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	reservation, err := c.service.GetReservationByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Reservation not found",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// CreateReservationHandler godoc
// @Summary Create a new reservation
// @Description Create a new reservation with the provided details
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body ReservationRequest true "Reservation details"
// @Success 201 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reservations [post]
func (c *ReservationController) CreateReservationHandler(ctx *gin.Context) {
	var req ReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	reservation, err := c.service.CreateReservation(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to create reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, reservation.ToResponse())
}

// UpdateReservationHandler godoc
// @Summary Update a reservation
// @Description Update reservation details by reservation ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Param reservation body ReservationRequest true "Updated reservation details"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reservations/{id} [put]
func (c *ReservationController) UpdateReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	var req ReservationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	reservation, err := c.service.UpdateReservation(id, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// DeleteReservationHandler godoc
// @Summary Delete a reservation
// @Description Delete reservation by reservation ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /reservations/{id} [delete]
func (c *ReservationController) DeleteReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	err = c.service.DeleteReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to delete reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// UpdateReservationStatusHandler godoc
// @Summary Update reservation status
// @Description Update the status of a reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Param status body StatusUpdateRequest true "New status"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reservations/{id}/status [patch]
func (c *ReservationController) UpdateReservationStatusHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	var req StatusUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	reservation, err := c.service.UpdateReservationStatus(id, req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to update status",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// CancelReservationHandler godoc
// @Summary Cancel a reservation
// @Description Cancel a reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reservations/{id}/cancel [post]
func (c *ReservationController) CancelReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	reservation, err := c.service.CancelReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to cancel reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// ConfirmReservationHandler godoc
// @Summary Confirm a reservation
// @Description Confirm a pending reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reservations/{id}/confirm [post]
func (c *ReservationController) ConfirmReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	reservation, err := c.service.ConfirmReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to confirm reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// CheckInReservationHandler godoc
// @Summary Check-in a reservation
// @Description Mark a reservation as checked in
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reservations/{id}/checkin [post]
func (c *ReservationController) CheckInReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	reservation, err := c.service.CheckInReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to check-in reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}

// CheckOutReservationHandler godoc
// @Summary Check-out a reservation
// @Description Mark a reservation as checked out
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path string true "Reservation ID (UUID)"
// @Success 200 {object} ReservationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /reservations/{id}/checkout [post]
func (c *ReservationController) CheckOutReservationHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	reservation, err := c.service.CheckOutReservation(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to check-out reservation",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation.ToResponse())
}
