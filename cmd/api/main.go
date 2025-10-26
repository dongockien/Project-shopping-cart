package main

import (
	"path/filepath"

	"dnk.com/hoc-golang/internal/app"
	"dnk.com/hoc-golang/internal/config"
	"dnk.com/hoc-golang/internal/utils"
	"dnk.com/hoc-golang/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {

	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")
	
	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	})

	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {

		logger.Log.Warn().Msg("No .env file found")
	} else {

		logger.Log.Info().Msg("Loaded successfully .env in api proccess")
	}

	// Inittialize configuration
	cfg := config.NewConfig()

	// Inittialize application
	application,err  := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}
	// Start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run failed")
	}
}

