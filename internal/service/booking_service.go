package service

import "github.com/zhxauda9/StayMate/internal/dal"

type Booking_serv_Impl interface{}

type booking_service struct {
	booking_repo dal.Booking_repo_Impl
}

func NewBookingService(booking_repo dal.Booking_repo_Impl) *booking_service {
	return &booking_service{booking_repo: booking_repo}
}
