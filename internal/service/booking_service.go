package service

import (
	"errors"
	"fmt"

	"github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/models"
)

type bookingService struct {
	bookingRepo postgres.BookingRepo
}

type BookingServ interface {
	CreateBooking(booking models.Booking) error
	GetBookingByID(id int) (models.Booking, error)
	GetAllBookings(sort, filterStart, filterEnd string, page int) ([]models.Booking, error)
	UpdateBooking(id int, booking models.Booking) error
	DeleteBooking(id int) error
}

func NewBookingService(bookingRepo postgres.BookingRepo) BookingServ {
	return &bookingService{bookingRepo: bookingRepo}
}

func (s *bookingService) CreateBooking(booking models.Booking) error {
	if !s.bookingRepo.CheckUserExists(booking.UserID) {
		return errors.New("user does not exist")
	}

	if !s.bookingRepo.CheckRoomExists(booking.RoomID) {
		return errors.New("room does not exist")
	}

	if s.bookingRepo.BookingExists(booking.RoomID, booking.CheckIn, booking.CheckOut) {
		return errors.New("room already booked for the selected dates")
	}

	return s.bookingRepo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id int) (models.Booking, error) {
	booking, err := s.bookingRepo.GetBookingByID(id)
	if err != nil {
		return models.Booking{}, fmt.Errorf("error in service layer while fetching booking by ID: %v", err)
	}
	return booking, nil
}

func (s *bookingService) UpdateBooking(id int, booking models.Booking) error {
	err := s.bookingRepo.UpdateBooking(id, booking)
	if err != nil {
		return fmt.Errorf("error in service layer while updating booking: %v", err)
	}
	return nil
}

func (s *bookingService) DeleteBooking(id int) error {
	err := s.bookingRepo.DeleteBooking(id)
	if err != nil {
		return fmt.Errorf("error in service layer while deleting booking: %v", err)
	}
	return nil
}

func (s *bookingService) GetAllBookings(sort, filterStart, filterEnd string, page int) ([]models.Booking, error) {
	const limit = 10
	offset := (page - 1) * limit

	bookings, err := s.bookingRepo.GetAllBookings(sort, filterStart, filterEnd, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching filtered bookings: %v", err)
	}
	return bookings, nil
}
