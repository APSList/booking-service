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
	val, exists := ctx.Get("organization_id")

	if !exists {
		ctx.JSON(401, gin.H{"error": "Organization ID not found in session"})
		return
	}

	// 2. Type-assert the value to int64
	orgID, ok := val.(int64)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Internal server error: ID format mismatch"})
		return
	}

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

	var filteredCustomers []*pb.Customer
	for _, cust := range resp.Customers {
		if cust.OrganizationId == orgID {
			filteredCustomers = append(filteredCustomers, cust)
		}
	}

	ctx.JSON(http.StatusOK, resp.Customers)
}

// CreateCustomerHandler calls CreateCustomer RPC
func (c *CustomerController) CreateCustomerHandler(ctx *gin.Context) {
	var req struct {
		FullName       string `json:"full_name"`
		Email          string `json:"email"`
		OrganizationId int64  `json:"organization_id"`
	}

	orgID, ok := c.getOrgID(ctx)
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := c.client.CreateCustomer(ctx, &pb.CreateCustomerRequest{
		FullName:       req.FullName,
		Email:          req.Email,
		OrganizationId: orgID,
	})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp.Customer)
}

// GetCustomerByIDHandler calls GetCustomer RPC
func (c *CustomerController) GetCustomerByIDHandler(ctx *gin.Context) {
	orgID, ok := c.getOrgID(ctx)
	if !ok {
		return
	}

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

	if resp.Customer.OrganizationId != orgID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this customer"})
		return
	}

	ctx.JSON(http.StatusOK, resp.Customer)
}

/*// UpdateCustomerHandler updates a customer (assuming you implement Update RPC)
func (c *CustomerController) UpdateCustomerHandler(ctx *gin.Context) {
	orgID, ok := c.getOrgID(ctx)
	if !ok {
		return
	}

	ctx.JSON(http.StatusNotImplemented, gin.H{"error": "UpdateCustomer RPC not implemented in proto"})
}*/

// DeleteCustomerHandler
func (c *CustomerController) DeleteCustomerHandler(ctx *gin.Context) {
	orgID, ok := c.getOrgID(ctx)
	if !ok {
		return
	}

	id, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	// STEP 1: Fetch the customer first to check ownership
	checkResp, err := c.client.GetCustomer(ctx, &pb.GetCustomerRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	// STEP 2: Verify OrgID
	if checkResp.Customer.OrganizationId != orgID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized delete attempt"})
		return
	}

	// STEP 3: Proceed with delete
	resp, err := c.client.DeleteCustomer(ctx, &pb.DeleteCustomerRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": resp.Success})
}

func (c *CustomerController) getOrgID(ctx *gin.Context) (int64, bool) {
	val, exists := ctx.Get("organization_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Organization ID not found in session"})
		return 0, false
	}
	orgID, ok := val.(int64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error: ID format mismatch"})
		return 0, false
	}
	return orgID, true
}
