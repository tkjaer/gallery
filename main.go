package main

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

var year = time.Now().Year()

func main() {
	// Get the log level from the environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info" // Default to info if not set
	}
	// Map the log level string to slog.Level
	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default to info if the input is invalid
	}

	// Get the add source option from the environment variable
	// This option determines whether to include the source file and line number in the logs
	addSourceEnv := os.Getenv("ADD_SOURCE")
	addSource := strings.ToLower(addSourceEnv) == "true"

	// Set up slog with the specified log level
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level, AddSource: addSource}))
	slog.SetDefault(logger)

	slog.Debug("Starting application", "timestamp", time.Now().Format(time.RFC3339))

	err := LoadConfig("config.yml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	err = process()
	if err != nil {
		slog.Error("Failed to process", "error", err)
		os.Exit(1)
	}

	slog.Debug("Application finished successfully", "timestamp", time.Now().Format(time.RFC3339))
}
