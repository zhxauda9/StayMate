package dal

import (
	"database/sql"
	"fmt"
	"github.com/zhxauda9/StayMate/models"
)

type BookingRepo interface {
	CreateBooking(booking models.Booking) error
	GetBookingByID(id int) (models.Booking, error)
	GetAllBookings() ([]models.Booking, error)
	UpdateBooking(id int, booking models.Booking) error
	DeleteBooking(id int) error
	CheckUserExists(userID int) (bool, error)
	CheckRoomExists(roomID int) (bool, error)
}

type bookingRepository struct {
	db *sql.DB
}

func (r *bookingRepository) CreateBooking(booking models.Booking) error {
	query := `
		INSERT INTO bookings (user_id, room_id, check_in, check_out) 
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := r.db.QueryRow(query, booking.UserID, booking.RoomID, booking.CheckIn, booking.CheckOut).Scan(&booking.ID)
	if err != nil {
		return fmt.Errorf("error inserting booking: %v", err)
	}
	return nil
}

func (r *bookingRepository) GetBookingByID(id int) (models.Booking, error) {
	query := `
		SELECT id, user_id, room_id, check_in, check_out
		FROM bookings WHERE id = $1
	`
	var booking models.Booking
	err := r.db.QueryRow(query, id).Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.CheckIn, &booking.CheckOut)
	if err != nil {
		return models.Booking{}, fmt.Errorf("error fetching booking by ID: %v", err)
	}
	return booking, nil
}

func (r *bookingRepository) GetAllBookings() ([]models.Booking, error) {
	query := `
		SELECT id, user_id, room_id, check_in, check_out 
		FROM bookings
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all bookings: %v", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.CheckIn, &booking.CheckOut)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *bookingRepository) UpdateBooking(id int, booking models.Booking) error {
	query := `
		UPDATE bookings
		SET user_id = $1, room_id = $2, check_in = $3, check_out = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(query, booking.UserID, booking.RoomID, booking.CheckIn, booking.CheckOut, id)
	if err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}
	return nil
}

func (r *bookingRepository) DeleteBooking(id int) error {
	query := `
		DELETE FROM bookings WHERE id = $1
	`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting booking: %v", err)
	}
	return nil
}

func NewBookingRepository(db *sql.DB) BookingRepo {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) CheckUserExists(userID int) (bool, error) {
	var count int
	query := "SELECT COUNT(1) FROM users WHERE id = $1"
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bookingRepository) CheckRoomExists(roomID int) (bool, error) {
	var count int
	query := "SELECT COUNT(1) FROM rooms WHERE id = $1"
	err := r.db.QueryRow(query, roomID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
