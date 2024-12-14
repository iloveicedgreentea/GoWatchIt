package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

// loggerKey is used as the key for storing the logger in the context
type loggerKey struct{}

var defaultLogger *slog.Logger

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func getLogFilePath() string {
	env := os.Getenv("LOG_ENV")
	if env == "" {
		return "/data/app.log"
	}

	return "./app.log"
}

// AddLoggerToContext adds a slog.Logger to the context
func AddLoggerToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLoggerFromContext retrieves the slog.Logger from the context
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	logger, ok := ctx.Value(loggerKey{}).(*slog.Logger)
	if !ok {
		// Return the default logger if not found in context
		return GetLogger()
	}
	return logger
}

// GetLogger returns the default slog.Logger instance
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		var handler slog.Handler

		// Determine the log level
		level := slog.LevelInfo
		if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
			level = slog.LevelDebug
		}

		// Set up logging to file if LOG_FILE is "true"
		if os.Getenv("LOG_FILE") == "true" {
			logFilePath := getLogFilePath()
			// Remove old log file
			err := os.Remove(logFilePath)
			if err != nil && !os.IsNotExist(err) {
				slog.Error("Failed to remove log file", "error", err)
			}

			// Open a new log file
			// #nosec G304 - We are not using user input to create the file
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
			if err != nil {
				slog.Error("Failed to open log file", "error", err)
			} else {
				// Create a multi-writer for both file and stdout
				multiWriter := io.MultiWriter(file, os.Stdout)
				handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
					Level:     level,
					AddSource: true,
				})
			}
		}

		// If no file handler was created, use a default stdout handler
		if handler == nil {
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			})
		}

		// Create the logger
		defaultLogger = slog.New(handler)

		// Set as default logger
		slog.SetDefault(defaultLogger)
	}

	return defaultLogger
}

func Error(msg string, args ...any) {
	log := GetLogger()
	log.Error(msg, args...)
}
