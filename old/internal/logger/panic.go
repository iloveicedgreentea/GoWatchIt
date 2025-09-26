package logger

import (
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
)

var logMu sync.Mutex

// PanicLogger wraps the standard logger to ensure panics are captured
func PanicLogger(next func()) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()

			// Ensure we have a logger even if panic happens during logger init
			logger := GetLogger()
			if logger == nil {
				// Emergency logging if logger isn't initialized
				handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: slog.LevelError,
				})
				logger = slog.New(handler)
			}

			logger.Error("PANIC RECOVERED",
				slog.Any("panic", err),
				slog.String("stack", string(stack)),
			)

			// Ensure logs are flushed
			if logFile != nil {
				_ = logFile.Sync()
			}

			// Re-panic to let the program crash with the original error
			panic(err)
		}
	}()
	next()
}

// WrapHandler wraps an HTTP handler with panic recovery
func WrapHandler(handler func()) func() {
	return func() {
		PanicLogger(handler)
	}
}
