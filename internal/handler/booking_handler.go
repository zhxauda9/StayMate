package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	l "github.com/zhxauda9/StayMate/internal/myLogger"
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
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to create a new booking.")

	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		l.Log.Error().Err(err).Msg("Error decoding booking data")
		http.Error(w, "Error decoding booking data", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(booking); err != nil {
		l.Log.Warn().Err(err).Msg("Validation error for booking data")
		http.Error(w, "Validation error for booking data", http.StatusBadRequest)
		return
	}

	err := h.bookingService.CreateBooking(booking)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error creating booking")
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Msg("Booking created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

// handler.go
func (h *bookingHandler) GetBookings(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch all bookings.")

	filterStart := r.URL.Query().Get("filterStart")
	filterEnd := r.URL.Query().Get("filterEnd")
	sort := r.URL.Query().Get("sort")
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	bookings, err := h.bookingService.GetAllBookings(sort, filterStart, filterEnd, page)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error fetching bookings")
		http.Error(w, "Error fetching bookings", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Msg("Fetched filtered bookings successfully")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}

func (h *bookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to fetch a specific booking.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid booking ID")
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.GetBookingByID(id)
	if err != nil {
		l.Log.Error().Err(err).Int("BookingID", id).Msg("Error fetching booking")
		http.Error(w, "Error fetching booking", http.StatusNotFound)
		return
	}
	l.Log.Info().Int("BookingID", id).Msg("Fetched booking successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) PutBooking(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to update a booking.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid booking ID")
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}
	l.Log.Debug().Int("BookingID", id).Msg("Parsed booking ID successfully")

	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		l.Log.Error().Err(err).Msg("Error decoding booking data")
		http.Error(w, "Error decoding booking data", http.StatusBadRequest)
		return
	}

	err = h.bookingService.UpdateBooking(id, booking)
	if err != nil {
		l.Log.Error().Err(err).Int("BookingID", id).Msg("Error updating booking")
		http.Error(w, "Error updating booking", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("BookingID", id).Msg("Updated booking successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func (h *bookingHandler) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	l.Log.Info().Str("IP", r.RemoteAddr).Msg("Received request to delete a booking.")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Log.Warn().Str("ID", idStr).Err(err).Msg("Invalid booking ID")
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	err = h.bookingService.DeleteBooking(id)
	if err != nil {
		l.Log.Error().Err(err).Int("BookingID", id).Msg("Error deleting booking")
		http.Error(w, "Error deleting booking", http.StatusInternalServerError)
		return
	}
	l.Log.Info().Int("BookingID", id).Msg("Deleted booking successfully")

	w.WriteHeader(http.StatusNoContent)
}
