package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func GetReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		db: db,
	}
}

// GetReservations returns all reservations ordered by creation date
func (r *ReservationRepository) GetReservations() ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationByID returns a single reservation by ID
func (r *ReservationRepository) GetReservationByID(id int) (*Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE id = $1
    `

	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservation, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *ReservationRepository) CreateReservation(reservation *Reservation) (*Reservation, error) {
	query := `
        INSERT INTO reservation (
            id, organization_id, property_id, customer_id, check_in_date, check_out_date,
            status, total_price, payment_url, price_elements, no_of_guests, 
            guest_data, additional_requests, created_at, update_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
        RETURNING id, organization_id, property_id, customer_id, check_in_date, check_out_date,
                  status, total_price, payment_url, price_elements, no_of_guests, 
                  guest_data, additional_requests, created_at, update_at
    `

	// Create a new instance to scan into
	res := &Reservation{}

	err := r.db.QueryRow(context.Background(), query,
		reservation.ID,
		reservation.OrganizationID,
		reservation.PropertyID,
		reservation.CustomerID,
		reservation.CheckInDate,
		reservation.CheckOutDate,
		reservation.Status,
		reservation.TotalPrice,
		reservation.PaymentURL, // $9
		reservation.PriceElements,
		reservation.NoOfGuests,
		reservation.GuestData,
		reservation.AdditionalRequests,
		reservation.CreatedAt,
		reservation.UpdatedAt, // $15
	).Scan(
		&res.ID,
		&res.OrganizationID,
		&res.PropertyID,
		&res.CustomerID,
		&res.CheckInDate,
		&res.CheckOutDate,
		&res.Status,
		&res.TotalPrice,
		&res.PaymentURL,
		&res.PriceElements,
		&res.NoOfGuests,
		&res.GuestData,
		&res.AdditionalRequests,
		&res.CreatedAt,
		&res.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to insert reservation: %w", err)
	}

	return res, nil
}

func (r *ReservationRepository) UpdateReservation(reservation *Reservation) (*Reservation, error) {
	query := `
        UPDATE reservation
        SET organization_id = $2,
            property_id = $3,
            customer_id = $4,
            check_in_date = $5,
            status = $6,
            total_price = $7,
            payment_url = $8,          -- Added this
            price_elements = $9,
            no_of_guests = $10,
            guest_data = $11,
            additional_requests = $12,
            check_out_date = $13,
            update_at = $14           -- Changed update_at to updated_at
        WHERE id = $1
        RETURNING id, organization_id, property_id, customer_id, check_in_date, 
                  check_out_date, status, total_price, payment_url, price_elements, 
                  no_of_guests, guest_data, additional_requests, created_at, update_at
    `

	rows, err := r.db.Query(
		context.Background(),
		query,
		reservation.ID,                 // $1
		reservation.OrganizationID,     // $2
		reservation.PropertyID,         // $3
		reservation.CustomerID,         // $4
		reservation.CheckInDate,        // $5
		reservation.Status,             // $6
		reservation.TotalPrice,         // $7
		reservation.PaymentURL,         // $8 - New
		reservation.PriceElements,      // $9
		reservation.NoOfGuests,         // $10
		reservation.GuestData,          // $11
		reservation.AdditionalRequests, // $12
		reservation.CheckOutDate,       // $13
		reservation.UpdatedAt,          // $14
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// RowToStructByName matches struct tags (db:"payment_url") to SQL column names
	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}

	return &updated, nil
}

// DeleteReservation deletes a reservation by ID
func (r *ReservationRepository) DeleteReservation(id int) error {
	query := `DELETE FROM reservation WHERE id = $1`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("reservation not found")
	}

	return nil
}

// CheckPropertyAvailability checks if a property is available for the given dates
func (r *ReservationRepository) CheckPropertyAvailability(propertyID int, checkIn, checkOut time.Time) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM reservation
            WHERE property_id = $1
              AND status NOT IN ('cancelled', 'rejected')
              AND (
                (check_in_date <= $2 AND check_out_date > $2) OR
                (check_in_date < $3 AND check_out_date >= $3) OR
                (check_in_date >= $2 AND check_out_date <= $3)
              )
        )
    `

	var exists bool
	err := r.db.QueryRow(context.Background(), query, propertyID, checkIn, checkOut).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// CheckPropertyAvailabilityExcluding checks availability excluding a specific reservation
func (r *ReservationRepository) CheckPropertyAvailabilityExcluding(excludeID int, propertyID int, checkIn, checkOut time.Time) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1
            FROM reservation
            WHERE property_id = $1
              AND id != $2
              AND status NOT IN ('cancelled', 'rejected')
              AND (
                (check_in_date <= $3 AND check_out_date > $3) OR
                (check_in_date < $4 AND check_out_date >= $4) OR
                (check_in_date >= $3 AND check_out_date <= $4)
              )
        )
    `

	var exists bool
	err := r.db.QueryRow(context.Background(), query, propertyID, excludeID, checkIn, checkOut).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetReservationsByCustomer returns all reservations for a customer
func (r *ReservationRepository) GetReservationsByCustomer(customerID int) ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE customer_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationsByProperty returns all reservations for a property
func (r *ReservationRepository) GetReservationsByProperty(propertyID int) ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE property_id = $1
        ORDER BY check_in_date DESC
    `

	rows, err := r.db.Query(context.Background(), query, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationsByOrganization returns all reservations for an organization
func (r *ReservationRepository) GetReservationsByOrganization(organizationID int) ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE organization_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationsByStatus returns reservations with a specific status
func (r *ReservationRepository) GetReservationsByStatus(status string) ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE status = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetUpcomingReservations returns future reservations
func (r *ReservationRepository) GetUpcomingReservations() ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE check_in_date > NOW()
          AND status NOT IN ('cancelled', 'rejected', 'completed')
        ORDER BY check_in_date ASC
    `

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

// GetReservationsByDateRange returns reservations within a date range
func (r *ReservationRepository) GetReservationsByDateRange(startDate, endDate time.Time) ([]Reservation, error) {
	query := `
        SELECT id, organization_id, property_id, customer_id, check_in_date, status, 
               total_price, price_elements, no_of_guests, guest_data, additional_requests, 
               check_out_date, created_at, update_at
        FROM reservation
        WHERE check_in_date >= $1 AND check_in_date <= $2
        ORDER BY check_in_date ASC
    `

	rows, err := r.db.Query(context.Background(), query, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations, err := pgx.CollectRows(rows, pgx.RowToStructByName[Reservation])
	if err != nil {
		return nil, err
	}

	return reservations, nil
}
