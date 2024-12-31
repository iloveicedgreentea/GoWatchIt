package logger

import (
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"sync"
)

var logMu sync.Mutex

// PanicLogger wraps the standard logger to ensure panics are captured
func PanicLogger(next func()) {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			logMu.Lock()
			defer logMu.Unlock()

			// Ensure we have a logger
			log := GetLogger()

			// Log both the panic and the stack trace
			log.Error("PANIC RECOVERED",
				slog.Any("error", err),
				slog.String("stack", string(stack)),
			)

			// Re-panic if we're in development mode
			if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
				panic(err)
			}

			// In production, exit with error code
			os.Exit(1)
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
