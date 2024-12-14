package handler

import (
	"net/http"

	"github.com/zhxauda9/StayMate/internal/service"
)

type bookingHandler struct {
	booking_service service.Booking_serv_Impl
}

func NewBookingHandler(booking_service service.Booking_serv_Impl) *bookingHandler {
	return &bookingHandler{booking_service: booking_service}
}

func (h *bookingHandler) PostBooking(w http.ResponseWriter, r *http.Request) {
}

func (h *bookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
}

func (h *bookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
}

func (h *bookingHandler) PutBooking(w http.ResponseWriter, r *http.Request) {
}

func (h *bookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
}
