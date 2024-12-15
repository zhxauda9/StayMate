package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type bookingHandler struct {
	booking_service service.BookingServ
}

func (h *bookingHandler) PostBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding booking data: %v", err), http.StatusBadRequest)
		return
	}

	err := h.booking_service.CreateBooking(booking)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating booking: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.booking_service.GetAllBookings()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching bookings: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}

func (h *bookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting id: %v", err), http.StatusBadRequest)
		return
	}

	booking, err := h.booking_service.GetBookingByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching booking: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) PutBooking(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting id: %v", err), http.StatusBadRequest)
		return
	}
	var booking models.Booking
	// Парсим данные для обновления из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding booking data: %v", err), http.StatusBadRequest)
		return
	}
	err = h.booking_service.UpdateBooking(id, booking)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating booking: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting id: %v", err), http.StatusBadRequest)
		return
	}

	err = h.booking_service.DeleteBooking(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting booking: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func NewBookingHandler(booking_service service.BookingServ) *bookingHandler {
	return &bookingHandler{booking_service: booking_service}
}
