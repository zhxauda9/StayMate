package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type bookingHandler struct {
	bookingService service.BookingServ
	validate       *validator.Validate
}

func NewBookingHandler(bookingService service.BookingServ) *bookingHandler {
	return &bookingHandler{
		bookingService: bookingService,
		validate:       validator.New(),
	}
}

func (h *bookingHandler) PostBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding booking data: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(booking); err != nil {
		http.Error(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	err := h.bookingService.CreateBooking(booking)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating booking: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := h.bookingService.GetAllBookings()
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

	booking, err := h.bookingService.GetBookingByID(id)
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
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding booking data: %v", err), http.StatusBadRequest)
		return
	}
	err = h.bookingService.UpdateBooking(id, booking)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating booking: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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

	err = h.bookingService.DeleteBooking(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting booking: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
