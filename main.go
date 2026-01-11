package main

import (
	_ "hostflow/booking-service/docs" // Import generated swagger docs
	"hostflow/booking-service/internal/bootstrap"
	"hostflow/booking-service/internal/communication"
	"hostflow/booking-service/internal/customer"
	"hostflow/booking-service/internal/kafka"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

// @title Hostflow Booking Service API
// @version 1.0
// @description This is a comprehensive booking/reservation service API for managing property reservations with full CRUD operations, status management, and availability checking.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@hostflow.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host hostflow.software/booking
// @BasePath /
// @schemes https

// @tag.name reservations
// @tag.description Operations related to reservations

// @tag.name customers
// @tag.description Customer-specific reservation operations

func main() {
	_ = godotenv.Load()

	fx.New(
		bootstrap.Module,
		kafka.Module,
		customer.Module,
		communication.Module,
	).Run()
}
