package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhxauda9/StayMate/internal/config"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/server"
)

const address = ":8080"

func InitApp() {
	// Load environment variables and initialize logger
	config.LoadEnvVariables()
	l.Log = l.NewZeroLoggerV2()
	if l.Log == nil {
		return
	}
	l.Log.Info().Msg("Starting the application. Trying to initialize server...")

	mux, err := server.InitServer()
	if err != nil {
		l.Log.Fatal().Err(err).Msg("Error initializing server")
	}

	// Create an HTTP server instance
	srv := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	// Channel to capture system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Start the server in a goroutine
	go func() {
		l.Log.Info().Str("Address", address).Msg("Server initialized successfully, starting HTTP server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Log.Fatal().Err(err).Msg("Error starting HTTP server")
		}
	}()

	// Wait for shutdown signal. The code won't go further till signal received
	<-quit
	l.Log.Warn().Msg("Shutting down server...")

	// Create a context with a timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		l.Log.Error().Err(err).Msg("Server forced to shutdown due to error")
	} else {
		l.Log.Info().Msg("Server shutdown gracefully")
	}

	l.Log.Info().Msg("Application exited successfully")
}
