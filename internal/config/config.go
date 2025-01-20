package config

import (
	"github.com/joho/godotenv"
	"github.com/zhxauda9/StayMate/internal/myLogger"
)

var AvailiableMimeTypes = map[string]bool{
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		myLogger.Log.Warn().Err(err).Msg("Could not load environment variables")
	}
}
