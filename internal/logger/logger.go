package logger

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func GetLogger() *log.Logger {
	logger := log.New()

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

	return logger
}
