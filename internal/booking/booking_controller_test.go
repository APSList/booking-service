package booking

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReservationService struct {
	mock.Mock
}

func (m *MockReservationService) GetReservations(orgID int64) ([]Reservation, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) CreateReservation(req *ReservationRequest, orgID int64) (*Reservation, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) UpdateReservation(id int, req *ReservationRequest, orgID int64) (*Reservation, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) DeleteReservation(id int, orgID int64) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) UpdateReservationStatus(id int, status string, orgID int64) (*Reservation, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) CancelReservation(id int, orgID int64) (*Reservation, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) ConfirmPayment(reservationID int) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockReservationService) GetReservationByID(id int, orgID int64) (*Reservation, error) {
	args := m.Called(id, orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Reservation), args.Error(1)
}

// TEST 1: Preverjanje zdravja (Liveness)
func TestLivenessHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	routes := ReservationRoutes{}
	r.GET("/health/liveness", routes.LivenessHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/liveness", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
}

// TEST 2: Uspešna pridobitev rezervacije po ID-ju
func TestGetReservationByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockReservationService)
	controller := GetReservationController(mockSvc)

	r := gin.Default()
	r.GET("/reservations/:id", func(c *gin.Context) {
		c.Set("organization_id", int64(100)) // Simulacija auth middleware-a
		controller.GetReservationByIDHandler(c)
	})

	// Nastavimo pričakovanje: ID=5, OrgID=100
	mockReservation := &Reservation{ID: 5, OrganizationID: 100, Status: "CONFIRMED"}
	mockSvc.On("GetReservationByID", 5, int64(100)).Return(mockReservation, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reservations/5", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":5`)
	mockSvc.AssertExpectations(t)
}

// TEST 3: Napaka - Manjkajoč Organization ID (Unauthorized)
func TestGetReservations_NoOrgID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockReservationService)
	controller := GetReservationController(mockSvc)

	r := gin.Default()
	r.GET("/reservations", controller.GetReservationsHandler)
	// Tu namerno NE nastavimo c.Set("organization_id", ...)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/reservations", nil)
	r.ServeHTTP(w, req)

	// Pričakujemo 401 Unauthorized, ker getOrgID() ne najde vrednosti
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Organization ID not found")
}
