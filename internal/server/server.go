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
	booking_repo := dal.NewBookingRepository(db)
	booking_service := service.NewBookingService(booking_repo)
	booking_handler := handler.NewBookingHandler(booking_service)

	mux.HandleFunc("POST /bookings", booking_handler.PostBooking)
	mux.HandleFunc("GET /bookings", booking_handler.GetBookings)
	mux.HandleFunc("GET /bookings/{id}", booking_handler.GetBooking)
	mux.HandleFunc("PUT /bookings/{id}", booking_handler.PutBooking)
	mux.HandleFunc("DELETE /bookings/{id}", booking_handler.DeleteBooking)

	user_repo := dal.NewUserRepository(db)
	user_service := service.NewUserService(user_repo)
	user_handler := handler.NewUserHandler(user_service)

	mux.HandleFunc("POST /user", user_handler.PostUser)
	mux.HandleFunc("GET /user", user_handler.GetUsers)
	mux.HandleFunc("GET /user/{id}", user_handler.GetUser)
	mux.HandleFunc("PUT /user/{id}", user_handler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", user_handler.DeleteUser)

	return mux, nil
}

func Connect_DB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error opening the database -> %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging the database -> %v", err)
	}
	return db, nil
}
