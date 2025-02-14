package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"payment-service/Mlogger"
	handlers "payment-service/handler"
	"payment-service/repositories"
	"payment-service/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment variables")
	}

	// Connect to Postgres
	pgPool, err := repositories.InitPostgresDB(os.Getenv("MICRO_PG_DSN"))
	if err != nil {
		log.Fatalf("Failed to connect Postgres DB: %v", err)
	}
	defer pgPool.Close()

	// Create the table if needed:
	ensureTransactionsTable(pgPool)

	r := gin.Default()
	r.Use(gin.Recovery()) // Handles panics
	r.Use(cors.Default()) // Enable CORS with default settings
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"}, // Replace with your frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	paymentHandler := &handlers.PaymentHandler{
		PGPool: pgPool,
		EmailSettings: services.EmailSettings{
			Host: os.Getenv("SMTP_HOST"),
			Port: os.Getenv("SMTP_PORT"),
			User: os.Getenv("EMAIL"),
			Pass: os.Getenv("PASSWORD"),
		},
		Logger: Mlogger.NewZeroLogger(),
	}

	// Payment route
	r.POST("/payment", paymentHandler.ProcessPayment)
	r.GET("/payment-history", paymentHandler.GetPaymentHistory)
	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + os.Getenv("MICRO_PORT"),
		Handler: r,
	}

	go func() {
		log.Printf("Payment service is  running on port %s", os.Getenv("MICRO_PORT"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Microservice server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down microservice server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Microservice server forced to shutdown: %v", err)
	}

	log.Println("Microservice server exited")
}

// ensureTransactionsTable can create the table if it doesn't exist
func ensureTransactionsTable(pgPool *pgxpool.Pool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sql := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			user_email TEXT NOT NULL,
			amount NUMERIC(10,2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			payment_method TEXT,
			card_details TEXT
		);
	`
	_, err := pgPool.Exec(ctx, sql)
	if err != nil {
		log.Printf("Failed to create transactions table: %v", err)
	}
}
