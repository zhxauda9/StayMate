package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/zhxauda9/StayMate/internal/dal"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/service"
)

func InitServer() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	db, err := Connect_DB()
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database -> %v", err)
	}

	// Booking
	booking_repo := dal.NewBookingRepository(db)
	booking_service := service.NewBookingService(booking_repo)
	booking_handler := handler.NewBookingHandler(booking_service)

	// Routing
	mux.HandleFunc("POST /bookings", booking_handler.PostBooking)
	mux.HandleFunc("GET /bookings", booking_handler.PostBooking)
	mux.HandleFunc("GET /bookings/{id}", booking_handler.PostBooking)
	mux.HandleFunc("PUT /bookings/{id}", booking_handler.PostBooking)
	mux.HandleFunc("DELETE /bookings/{id}", booking_handler.PostBooking)

	return mux, nil
}

func Connect_DB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening the database -> %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
