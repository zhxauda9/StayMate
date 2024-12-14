package service

import (
	"fmt"
	"github.com/zhxauda9/StayMate/internal/dal"
	"github.com/zhxauda9/StayMate/models"
)

type BookingServ interface {
	CreateBooking(booking models.Booking) error
	GetBookingByID(id int) (models.Booking, error)
	GetAllBookings() ([]models.Booking, error)
	UpdateBooking(id int, booking models.Booking) error
	DeleteBooking(id int) error
}

type bookingService struct {
	bookingRepo dal.BookingRepo
}

func (s *bookingService) CreateBooking(booking models.Booking) error {
	userExists, err := s.bookingRepo.CheckUserExists(booking.UserID)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}
	if !userExists {
		return fmt.Errorf("user with id %d does not exist", booking.UserID)
	}

	roomExists, err := s.bookingRepo.CheckRoomExists(booking.RoomID)
	if err != nil {
		return fmt.Errorf("error checking if room exists: %v", err)
	}
	if !roomExists {
		return fmt.Errorf("room with id %d does not exist", booking.RoomID)
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

func (s *bookingService) GetAllBookings() ([]models.Booking, error) {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return nil, fmt.Errorf("error in service layer while fetching all bookings: %v", err)
	}
	return bookings, nil
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

func NewBookingService(bookingRepo dal.BookingRepo) BookingServ {
	return &bookingService{bookingRepo: bookingRepo}
}
