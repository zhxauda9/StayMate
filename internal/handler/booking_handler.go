package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
	"github.com/zhxauda9/StayMate/models"
)

type bookingHandler struct {
	bookingService service.BookingServ
	roomService    service.RoomServ
	validate       *validator.Validate
}

func NewBookingHandler(bookingService service.BookingServ, roomService service.RoomServ) *bookingHandler {
	return &bookingHandler{
		bookingService: bookingService,
		roomService:    roomService,
		validate:       validator.New(),
	}
}

func (h *bookingHandler) PostBookingV2(w http.ResponseWriter, r *http.Request) {
	var request struct {
		UserID         int     `json:"user_id"`
		Email          string  `json:"email"`
		Address        string  `json:"address"`
		CardNumber     string  `json:"card_number"`
		ExpirationDate string  `json:"expiration_date"`
		CVV            string  `json:"cvv"`
		RoomID         int     `json:"room_id"`
		Price          float64 `json:"price"`
		CheckIn        string  `json:"check_in"`
		CheckOut       string  `json:"check_out"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		l.Log.Error().Err(err).Msg("Error decoding booking data")
		http.Error(w, "Error decoding booking data", http.StatusBadRequest)
		return
	}

	l.Log.Info().
		Int("User id", request.UserID).
		Str("Email", request.Email).
		Str("Address", request.Address).
		Str("Card Number", maskCardNumber(request.CardNumber)).
		Str("Expiration Date", request.ExpirationDate).
		Int("Room ID", request.RoomID).Msg("Received booking request")

	// Should send request to 8081 port localhost POST /payment
	//
	paymentRequest := map[string]interface{}{
		"user_id":         request.UserID,
		"email":           request.Email,
		"amount":          100.0,
		"card_number":     request.CardNumber,
		"expiration_date": request.ExpirationDate,
		"cvv":             request.CVV,
	}

	paymentResponse, err := sendPaymentRequest(paymentRequest)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error processing payment")
		http.Error(w, "Error processing payment", http.StatusInternalServerError)
		return
	}

	// Handle the response from the payment service
	if paymentResponse["status"] != "success" {
		l.Log.Error().Msg("Payment failed")
		http.Error(w, "Payment failed", http.StatusPaymentRequired)
		return
	}

	booking := models.Booking{
		UserID:   request.UserID,
		RoomID:   request.RoomID,
		CheckIn:  request.CheckIn,
		CheckOut: request.CheckOut,
	}

	err = h.bookingService.CreateBooking(booking)
	if err != nil {
		l.Log.Error().Err(err).Msg("Error creating booking")
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	err = h.roomService.SetStatus(booking.RoomID, "occupied")
	if err != nil {
		l.Log.Error().Err(err).Msg("Error creating booking")
		http.Error(w, "Error creating booking", http.StatusInternalServerError)
		return
	}

	l.Log.Info().Msg("Booking(V2) created and payment processed successfully")

	var responce = struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}{
		Message: "Booking and payment processed successfully",
		Status:  "Paid",
	}
	json.NewEncoder(w).Encode(responce)
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

func sendPaymentRequest(data map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payment request: %w", err)
	}

	url := "http://localhost:8081/payment"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to payment service: %w", err)
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding payment response: %w", err)
	}

	return response, nil
}

func maskCardNumber(cardNumber string) string {
	if len(cardNumber) <= 4 {
		return cardNumber
	}
	return strings.Repeat("*", len(cardNumber)-4) + cardNumber[len(cardNumber)-4:]
}
