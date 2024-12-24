package cmd

import (
	"net/http"

	"github.com/zhxauda9/StayMate/internal/config"
	l "github.com/zhxauda9/StayMate/internal/myLogger"
	"github.com/zhxauda9/StayMate/internal/server"
)

const address = "127.0.0.1:8080"

func InitApp() {
	config.LoadEnvVariables()
	l.Log = l.NewZeroLogger()
	l.Log.Info().Msg("Starting the application. Tryining to initialize server")

	mux, err := server.InitServer()
	if err != nil {
		l.Log.Fatal().Err(err).Msg("Error initializing server")
	}

	l.Log.Info().Msg("Server initialized successfully, starting HTTP server on address http://localhost:8080")
	if err := http.ListenAndServe(address, mux); err != nil {
		l.Log.Fatal().Err(err).Msg("Error starting HTTP server")
	}
}
