package dal

import (
	"fmt"

	"github.com/zhxauda9/StayMate/models"
	"gorm.io/gorm"
)

type BookingRepo interface {
	CreateBooking(booking models.Booking) error
	GetBookingByID(id int) (models.Booking, error)
	GetAllBookings() ([]models.Booking, error)
	UpdateBooking(id int, booking models.Booking) error
	DeleteBooking(id int) error
	CheckUserExists(userID int) bool
	BookingExists(roomID int, checkIn, checkOut string) bool
}

type bookingRepository struct {
	db *gorm.DB
}

func (r *bookingRepository) CreateBooking(booking models.Booking) error {
	if err := r.db.Create(&booking).Error; err != nil {
		return fmt.Errorf("error inserting booking: %v", err)
	}
	return nil
}

func (r *bookingRepository) GetBookingByID(id int) (models.Booking, error) {
	var booking models.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return models.Booking{}, fmt.Errorf("booking with ID %d not found", id)
		}
		return models.Booking{}, fmt.Errorf("error fetching booking by ID: %v", err)
	}
	return booking, nil
}

func (r *bookingRepository) GetAllBookings() ([]models.Booking, error) {
	var bookings []models.Booking
	if err := r.db.Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("error fetching all bookings: %v", err)
	}
	return bookings, nil
}

func (r *bookingRepository) CheckUserExists(userID int) bool {
	var count int64
	r.db.Model(&models.User{}).Where("id = ?", userID).Count(&count)
	return count > 0
}

func (r *bookingRepository) BookingExists(roomID int, checkIn, checkOut string) bool {
	var count int64
	r.db.Model(&models.Booking{}).
		Where("room_id = ? AND check_in < ? AND check_out > ?", roomID, checkOut, checkIn).
		Count(&count)
	return count > 0
}

func (r *bookingRepository) UpdateBooking(id int, booking models.Booking) error {
	var existingBooking models.Booking
	if err := r.db.First(&existingBooking, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("booking with ID %d not found", id)
		}
		return fmt.Errorf("error fetching booking for update: %v", err)
	}
	if err := r.db.Model(&existingBooking).Updates(booking).Error; err != nil {
		return fmt.Errorf("error updating booking: %v", err)
	}
	return nil
}

func (r *bookingRepository) DeleteBooking(id int) error {
	var booking models.Booking
	if err := r.db.First(&booking, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("booking with ID %d not found", id)
		}
		return fmt.Errorf("error fetching booking for deletion: %v", err)
	}
	if err := r.db.Delete(&booking).Error; err != nil {
		return fmt.Errorf("error deleting booking: %v", err)
	}
	return nil
}

func NewBookingRepository(db *gorm.DB) BookingRepo {
	return &bookingRepository{db: db}
}
