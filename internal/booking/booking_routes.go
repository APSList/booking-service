package booking

import (
	"hostflow/booking-service/pkg/lib"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ReservationRoutes struct
type ReservationRoutes struct {
	logger                lib.Logger
	router                *lib.Router
	reservationController *ReservationController
}

// SetReservationRoutes returns a ReservationRoutes struct
func SetReservationRoutes(
	logger lib.Logger,
	router *lib.Router,
	reservationController *ReservationController,
) ReservationRoutes {
	return ReservationRoutes{
		logger:                logger,
		router:                router,
		reservationController: reservationController,
	}
}

// Setup registers the reservation routes
func (route ReservationRoutes) Setup() {
	route.logger.Info("Setting up [RESERVATION] routes.")

	// Main reservations routes
	reservations := route.router.Group("/reservations")
	{
		reservations.GET("", route.reservationController.GetReservationsHandler)
		reservations.POST("", route.reservationController.CreateReservationHandler)
		reservations.GET("/:id", route.reservationController.GetReservationByIDHandler)
		reservations.PUT("/:id", route.reservationController.UpdateReservationHandler)
		reservations.DELETE("/:id", route.reservationController.DeleteReservationHandler)
	}

	// Swagger documentation route
	route.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	route.logger.Info("Swagger documentation available at: /swagger/index.html")
	route.logger.Info("[RESERVATION] routes setup complete.")
}
