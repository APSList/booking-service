package booking

import (
	"hostflow/booking-service/internal/communication"
	"hostflow/booking-service/internal/customer"
	"hostflow/booking-service/internal/middlewares"
	"hostflow/booking-service/pkg/lib"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ReservationRoutes struct
type ReservationRoutes struct {
	logger                  lib.Logger
	router                  *lib.Router
	reservationController   *ReservationController
	customerController      *customer.CustomerController
	communicationController *communication.CommunicationController
	authMiddleware          middlewares.AuthMiddleware
}

// SetReservationRoutes returns a ReservationRoutes struct
func SetReservationRoutes(
	logger lib.Logger,
	router *lib.Router,
	reservationController *ReservationController,
	customerController *customer.CustomerController,
	authMiddleware middlewares.AuthMiddleware,
	communicationController *communication.CommunicationController,
) ReservationRoutes {
	return ReservationRoutes{
		logger:                  logger,
		router:                  router,
		reservationController:   reservationController,
		customerController:      customerController,
		authMiddleware:          authMiddleware,
		communicationController: communicationController,
	}
}

// Setup registers the reservation routes
func (route ReservationRoutes) Setup() {
	route.logger.Info("Setting up [RESERVATION] routes.")

	// Main reservations routes
	reservations := route.router.Group("/reservations")
	reservations.Use(route.authMiddleware.Handler())
	{
		reservations.GET("", route.reservationController.GetReservationsHandler)
		reservations.POST("/", route.reservationController.CreateReservationHandler)
		reservations.GET("/:id", route.reservationController.GetReservationByIDHandler)
		reservations.PUT("/:id", route.reservationController.UpdateReservationHandler)
		reservations.DELETE("/:id", route.reservationController.DeleteReservationHandler)
	}

	customers := route.router.Group("/customer")
	customers.Use(route.authMiddleware.Handler())
	{
		customers.GET("", route.customerController.GetCustomerHandler)
		customers.POST("/", route.customerController.CreateCustomerHandler)
		customers.GET("/:id", route.customerController.GetCustomerByIDHandler)
		//customers.PUT("/:id", route.customerController.UpdateCustomerHandler)
		customers.DELETE("/:id", route.customerController.DeleteCustomerHandler)
	}

	communications := route.router.Group("/communication")
	{
		communications.POST("/email", route.communicationController.SendEmailHandler)
	}

	metrics := route.router.Group("/metrics")
	{
		metrics.GET("", gin.WrapH(promhttp.Handler()))
	}

	health := route.router.Group("/health")
	{
		health.GET("/live", route.LivenessHandler)
		health.GET("/ready", route.ReadinessHandler)
	}

	// Swagger documentation route
	route.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	route.logger.Info("Swagger documentation available at: /swagger/index.html")
	route.logger.Info("[RESERVATION] routes setup complete.")
}
