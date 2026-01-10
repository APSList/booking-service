package communication

import (
	"net/http"

	pb "hostflow/booking-service/internal/communication/proto"

	"github.com/gin-gonic/gin"
)

// CommunicationController wraps the gRPC client
type CommunicationController struct {
	client pb.CommunicationServiceClient
}

// NewController returns a new CommunicationController
func NewController(client pb.CommunicationServiceClient) *CommunicationController {
	return &CommunicationController{client: client}
}

// SendEmailHandler handles POST requests to trigger email sending via gRPC
func (c *CommunicationController) SendEmailHandler(ctx *gin.Context) {
	// 1. Define the expected JSON body from the frontend/caller
	var req struct {
		CustomerID int64  `json:"customer_id"`
		PropertyID string `json:"property_id"`
		Type       int32  `json:"type"` // 1: PAYMENT, 2: CONFIRMATION
		PaymentURL string `json:"payment_url"`
	}

	// 2. Bind incoming JSON to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 3. Call the gRPC SendEmail method
	resp, err := c.client.SendEmail(ctx, &pb.SendEmailRequest{
		CustomerId: req.CustomerID,
		PropertyId: req.PropertyID,
		Type:       pb.EmailType(req.Type),
		PaymentUrl: req.PaymentURL,
	})

	// 4. Handle gRPC communication errors (e.g., service down)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "Failed to communicate with email service: " + err.Error()})
		return
	}

	// 5. Return the result from the communication service
	if !resp.Success {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": resp.Message,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": resp.Message,
	})
}
