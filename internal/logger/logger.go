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

var (
	defaultLogger *slog.Logger
	logFile       *os.File // Keep track of the log file
)

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}

func InitLoggerFile() error {
	if os.Getenv("LOG_FILE") == "true" {
		logFilePath := getLogFilePath()

		// Remove old log file if it exists
		err := os.Remove(logFilePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		// Create new log file
		// #nosec G304 - We are not using user input to create the file
		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
		if err != nil {
			return err
		}

		// Store the file handle globally
		logFile = file

		// Create a new logger with the file
		handler := slog.NewJSONHandler(io.MultiWriter(file, os.Stdout), &slog.HandlerOptions{
			Level:     getLogLevel(),
			AddSource: true,
		})

		// Set up the default logger
		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	}
	return nil
}

func getLogLevel() slog.Level {
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func getLogFilePath() string {
	// if this is unset assume running in docker
	env, ok := os.LookupEnv("LOG_ENV")
	if env == "" || !ok {
		return "/data/app.log"
	}

	return "./applog.log"
}

// CleanupLogger ensures the log file is properly closed
func CleanupLogger() error {
	if logFile != nil {
		return logFile.Close()
	}
	return nil
}

// AddLoggerToContext adds a slog.Logger to the context
func AddLoggerToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// GetLoggerFromContext retrieves the slog.Logger from the context
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return logger
	}
	return GetLogger()
}

// GetLogger returns the default slog.Logger instance

func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		var handler slog.Handler
		level := getLogLevel()

		if os.Getenv("LOG_FILE") == "true" && logFile != nil {
			// Use existing log file if available
			multiWriter := io.MultiWriter(logFile, os.Stdout)
			handler = slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			})
		} else {
			// Fallback to stdout only
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     level,
				AddSource: true,
			})
		}

		defaultLogger = slog.New(handler)
		slog.SetDefault(defaultLogger)
	}
	return defaultLogger
}

func Error(msg string, args ...any) {
	log := GetLogger()
	log.Error(msg, args...)
}
