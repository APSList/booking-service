package customer

import (
	_ "context"
	pb "hostflow/booking-service/internal/customer/proto"
	"net/http"
	"strconv"

	//pb "github.com/hostflow/extra/customerpb"
	"github.com/gin-gonic/gin"
)

// CustomerController Controller wraps the gRPC client
type CustomerController struct {
	client pb.CustomerServiceClient
}

// NewController returns a CustomerController
func NewController(client pb.CustomerServiceClient) *CustomerController {
	return &CustomerController{client: client}
}

// GetCustomerHandler returns a list of customers (ListCustomers RPC)
func (c *CustomerController) GetCustomerHandler(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	resp, err := c.client.ListCustomers(ctx, &pb.ListCustomersRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp.Customers)
}

// CreateCustomerHandler calls CreateCustomer RPC
func (c *CustomerController) CreateCustomerHandler(ctx *gin.Context) {
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.client.CreateCustomer(ctx, &pb.CreateCustomerRequest{
		FullName: req.FullName,
		Email:    req.Email,
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp.Customer)
}

// GetCustomerByIDHandler calls GetCustomer RPC
func (c *CustomerController) GetCustomerByIDHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	resp, err := c.client.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, resp.Customer)
}

// UpdateCustomerHandler updates a customer (assuming you implement Update RPC)
func (c *CustomerController) UpdateCustomerHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "UpdateCustomer RPC not implemented in proto"})
}

// DeleteCustomerHandler calls DeleteCustomer RPC
func (c *CustomerController) DeleteCustomerHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer id"})
		return
	}

	resp, err := c.client.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": resp.Success})
}
