package logger

import (
	"io"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	logger *log.Logger
	once   sync.Once
)

func GetLogger() *log.Logger {
	once.Do(func() {
		logger = log.New()

		// log to file if LOG_FILE is not false
		if os.Getenv("LOG_FILE") != "false" {
			// truncate old logs
			err := os.Remove("/data/application.log")
			if err != nil && !os.IsNotExist(err) {
				log.Fatalf("Failed to remove log file: %v", err)
			}
			// Open a file for logging
			file, err := os.OpenFile("/data/application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("Failed to open log file: %v", err)
			}

			// Set logger output to the file and stdout
			logger.SetOutput(io.MultiWriter(file, os.Stdout))
		}
		// log level
		if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" {
			logger.SetLevel(log.DebugLevel)
		}
		if os.Getenv("SUPER_DEBUG") == "true" {
			logger.SetReportCaller(true)
		}

		logger.SetFormatter(&log.TextFormatter{
			TimestampFormat: "01-02-2006 15:04:05", FullTimestamp: true,
		})
	})
	return logger
}
