package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	dalpostgres "github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/internal/handler"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
)

func serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("web", "home.html"))
}

func InitServer() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	l.Log.Info().Msg("Trying to connect to database...")
	db, err := Connect_DB()
	if err != nil {
		l.Log.Error().Msg("Failed to connect to database")
		return nil, fmt.Errorf("error connecting to the database -> %v", err)
	}
	l.Log.Info().Msg("Successfully connected to database")

	// Serving static files (html, css, js, ...)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", serveHTML)

	booking_repo := dalpostgres.NewBookingRepository(db)
	booking_service := service.NewBookingService(booking_repo)
	booking_handler := handler.NewBookingHandler(booking_service)

	mux.HandleFunc("POST /bookings", booking_handler.PostBooking)
	mux.HandleFunc("GET /bookings", booking_handler.GetBookings)
	mux.HandleFunc("GET /bookings/{id}", booking_handler.GetBooking)
	mux.HandleFunc("PUT /bookings/{id}", booking_handler.PutBooking)
	mux.HandleFunc("DELETE /bookings/{id}", booking_handler.DeleteBooking)

	user_repo := dalpostgres.NewUserRepository(db)
	user_service := service.NewUserService(user_repo)
	user_handler := handler.NewUserHandler(user_service)

	mux.HandleFunc("POST /users", user_handler.CreateUser)
	mux.HandleFunc("GET /users", user_handler.GetUsers)
	mux.HandleFunc("GET /users/{id}", user_handler.GetUserByID)
	mux.HandleFunc("PUT /users/{id}", user_handler.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", user_handler.DeleteUser)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	mailServ, err := service.NewMailService(smtpHost, smtpPort, email, password)
	if err != nil {
		os.Exit(1)
	}
	mail_handler := handler.NewMailHandler(mailServ, user_service)

	mux.HandleFunc("GET /send-email", mail_handler.ServeMail)

	return mux, nil
}

func Connect_DB() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf(
		"user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening the database -> %v", err)
	}
	return db, nil
}
