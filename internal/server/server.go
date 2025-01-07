package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	dalpostgres "github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/middleware"
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

	// Serves home page
	mux.HandleFunc("/", serveHTML)

	rateLimitter := rate.NewLimiter(1, 3)                                  // Rate limit of 1 request per second with a burst of 3 requests
	limiterMiddlware := middleware.RateLimiterMiddlewareFunc(rateLimitter) // Middleware Function to wrap handlers
	// Init booking service
	booking_repo := dalpostgres.NewBookingRepository(db)
	booking_service := service.NewBookingService(booking_repo)
	booking_handler := handler.NewBookingHandler(booking_service)

	mux.Handle("POST /bookings", limiterMiddlware(http.HandlerFunc(booking_handler.PostBooking)))
	mux.Handle("GET /bookings", limiterMiddlware(http.HandlerFunc(booking_handler.GetBookings)))
	mux.Handle("GET /bookings/{id}", limiterMiddlware(http.HandlerFunc(booking_handler.GetBooking)))
	mux.Handle("PUT /bookings/{id}", limiterMiddlware(http.HandlerFunc(booking_handler.PutBooking)))
	mux.Handle("DELETE /bookings/{id}", limiterMiddlware(http.HandlerFunc(booking_handler.DeleteBooking)))

	// Init user service
	user_repo := dalpostgres.NewUserRepository(db)
	user_service := service.NewUserService(user_repo)
	user_handler := handler.NewUserHandler(user_service)

	mux.Handle("POST /users", limiterMiddlware(http.HandlerFunc(user_handler.CreateUser)))
	mux.Handle("GET /users", limiterMiddlware(http.HandlerFunc(user_handler.GetUsers)))
	mux.Handle("GET /users/{id}", limiterMiddlware(http.HandlerFunc(user_handler.GetUserByID)))
	mux.Handle("PUT /users/{id}", limiterMiddlware(http.HandlerFunc(user_handler.UpdateUser)))
	mux.Handle("DELETE /users/{id}", limiterMiddlware(http.HandlerFunc(user_handler.DeleteUser)))

	// Init Mailing service
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	mailServ, err := service.NewMailService(smtpHost, smtpPort, email, password)
	if err != nil {
		os.Exit(1)
	}
	mail_handler := handler.NewMailHandler(mailServ, user_service)

	mux.Handle("GET /mail", limiterMiddlware(http.HandlerFunc(mail_handler.ServeMail)))
	mux.Handle("POST /api/mail", limiterMiddlware(http.HandlerFunc(mail_handler.SendMailHandler)))
	mux.Handle("POST /api/mailFile", limiterMiddlware(http.HandlerFunc(mail_handler.SendMailFileHandler)))

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
