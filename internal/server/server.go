package server

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/zhxauda9/StayMate/internal/dal/migrations"
	dalpostgres "github.com/zhxauda9/StayMate/internal/dal/postgres"
	"github.com/zhxauda9/StayMate/internal/handler"
	"github.com/zhxauda9/StayMate/internal/middleware"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/service"
)

func InitServer() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	l.Log.Info().Msg("Trying to connect to database...")
	db, err := Connect_DB()
	if err != nil {
		l.Log.Error().Msg("Failed to connect to database")
		return nil, fmt.Errorf("error connecting to the database -> %v", err)
	}
	l.Log.Info().Msg("Successfully connected to database")

	// Migrations
	err = migrations.AutoMigrateDatabase(db)
	if err != nil {
		return nil, err
	}
	// Serving static files (html, css, js, ...)
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	rateLimitter := rate.NewLimiter(2, 2)                                 // Rate limit of 1 request per second with a burst of 3 requests
	limitMiddleware := middleware.RateLimiterMiddlewareFunc(rateLimitter) // Middleware Function to rate limit handlers
	logMiddleware := middleware.LoggingMiddlewareFunc(l.Log)              // Middleware Function for logging

	// Auth middlewares
	adminMid := middleware.MiddlewareWrapper(middleware.AuthAdminMiddleware)
	userMid := middleware.MiddlewareWrapper(middleware.AuthUserMiddleware)

	// Serves html pages
	mux.HandleFunc("/", handler.ServeHome)
	mux.HandleFunc("/register", handler.ServeRegister)
	mux.Handle("/login", http.HandlerFunc(handler.ServeLogin))
	mux.Handle("/profile", userMid(http.HandlerFunc(handler.ServeProfile)))
	mux.Handle("/verify-email", http.HandlerFunc(handler.ServeEmailVerify))

	mux.Handle("/admin", adminMid(http.HandlerFunc(handler.ServeAdmin)))
	mux.Handle("/admin/bookings", adminMid(http.HandlerFunc(handler.ServeBookings)))
	mux.Handle("/admin/rooms", adminMid(http.HandlerFunc(handler.ServeRooms)))
	mux.Handle("/admin/users", adminMid(http.HandlerFunc(handler.ServeUsers)))
	mux.Handle("/admin/chats", adminMid(http.HandlerFunc(handler.ServeAdminChats)))

	// Init booking service
	booking_repo := dalpostgres.NewBookingRepository(db)
	booking_service := service.NewBookingService(booking_repo)
	booking_handler := handler.NewBookingHandler(booking_service)

	mux.Handle("POST /bookings", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(booking_handler.PostBooking)))))
	mux.Handle("GET /bookings", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(booking_handler.GetBookings)))))
	mux.Handle("GET /bookings/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(booking_handler.GetBooking)))))
	mux.Handle("PUT /bookings/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(booking_handler.PutBooking)))))
	mux.Handle("DELETE /bookings/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(booking_handler.DeleteBooking)))))

	// Init user service
	user_repo := dalpostgres.NewUserRepository(db)
	user_service := service.NewUserService(user_repo)
	user_handler := handler.NewUserHandler(user_service)

	mux.Handle("POST /users", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(user_handler.CreateUser)))))
	mux.Handle("GET /users", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(user_handler.GetUsers)))))
	mux.Handle("GET /users/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(user_handler.GetUserByID)))))
	mux.Handle("PUT /users/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(user_handler.UpdateUser)))))
	mux.Handle("DELETE /users/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(user_handler.DeleteUser)))))

	room_repo := dalpostgres.NewRoomRepository(db)
	room_service := service.NewRoomService(room_repo)
	room_handler := handler.NewRoomHandler(room_service)

	mux.Handle("POST /rooms", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(room_handler.PostRoom)))))
	mux.Handle("GET /rooms", logMiddleware(limitMiddleware(http.HandlerFunc(room_handler.GetRooms))))
	mux.Handle("GET /rooms/{id}", logMiddleware(limitMiddleware(http.HandlerFunc(room_handler.GetRoom))))
	mux.Handle("PUT /rooms/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(room_handler.PutRoom)))))
	mux.Handle("DELETE /rooms/{id}", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(room_handler.DeleteRoom)))))

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")

	mailServ, err := service.NewMailService(smtpHost, smtpPort, email, password)
	if err != nil {
		os.Exit(1)
	}
	mail_handler := handler.NewMailHandler(mailServ, user_service)
	mux.Handle("GET /mail", logMiddleware(adminMid(limitMiddleware(http.HandlerFunc(mail_handler.ServeMail)))))
	mux.Handle("POST /api/mail", logMiddleware(limitMiddleware(http.HandlerFunc(mail_handler.SendMailHandler))))
	mux.Handle("POST /api/mailFile", logMiddleware(limitMiddleware(http.HandlerFunc(mail_handler.SendMailFileHandler))))

	verifyRepo := dalpostgres.NewVerifyRepository(db)
	authService := service.NewAuthService(user_service)
	authHandler := handler.NewAuthHandler(authService, verifyRepo)

	mux.Handle("POST /auth/register", logMiddleware(limitMiddleware(http.HandlerFunc(authHandler.Register))))
	mux.Handle("POST /auth/login", logMiddleware(limitMiddleware(http.HandlerFunc(authHandler.Login))))
	mux.Handle("GET /auth/validate", logMiddleware(limitMiddleware(http.HandlerFunc(authHandler.ValidateToken)))) // is it useful????
	mux.Handle("GET /auth/profile", logMiddleware(userMid(limitMiddleware(http.HandlerFunc(authHandler.GetProfile)))))
	mux.Handle("POST /auth/logout", logMiddleware(userMid(limitMiddleware((http.HandlerFunc(authHandler.Logout))))))

	verifyHandler := handler.NewVerifyHandler(db, mailServ)
	mux.HandleFunc("POST /auth/request-code", verifyHandler.SendVerifyCode)
	mux.HandleFunc("POST /auth/verify", verifyHandler.Verify)

	// Chat Handlers
	chatRepo := dalpostgres.NewChatRepository(db)
	chatHandler := handler.NewChatHandler(chatRepo, user_repo, l.Log)

	// Chats
	mux.Handle("POST /api/chat/start", logMiddleware(userMid(http.HandlerFunc(chatHandler.StartChat))))
	mux.Handle("GET /api/chat/history/{id}", logMiddleware(userMid(http.HandlerFunc(chatHandler.GetChatHistory))))
	mux.Handle("PUT /api/chat/close/{id}", logMiddleware(adminMid(http.HandlerFunc(chatHandler.CloseChat))))
	mux.Handle("GET /api/admin/chats", logMiddleware(adminMid(http.HandlerFunc(chatHandler.GetActiveChats))))
	// Serves template page for admin
	mux.Handle("GET /admin/chats/{id}", logMiddleware(adminMid(http.HandlerFunc(chatHandler.AdminChatPage))))

	// Websocket handlers for chat
	chatWebsocketHandler := handler.NewChatWebsocketHandler(l.Log, chatRepo)
	mux.Handle("/ws/user", http.HandlerFunc(chatWebsocketHandler.UserHandler))
	mux.Handle("/ws/admin", http.HandlerFunc(chatWebsocketHandler.AdminHandler))
	return mux, nil
}

func Connect_DB() (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	l.Log.Debug().Msg(psqlInfo)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error opening the database -> %v", err)
	}
	return db, nil
}
