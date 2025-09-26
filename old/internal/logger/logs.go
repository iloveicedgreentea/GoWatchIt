package logger

// sending logs to frontend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iloveicedgreentea/go-plex/models"
)

const MaxLogEntries = 1000 // Limit the number of log entries returned

func GetLogEntries() ([]models.LogEntry, error) {
	logFilePath := getLogFilePath()
	// #nosec G304 - We are not using user input to create the file
	file, err := os.Open(logFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Failed to close log file", err)
		}
	}()

	var entries []models.LogEntry
	scanner := bufio.NewScanner(file)

	// Create a larger buffer for scanner to handle long lines
	const maxCapacity = 512 * 1024 // 512KB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Read each line
	for scanner.Scan() {
		var entry models.LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			// Skip malformed lines
			return nil, fmt.Errorf("failed to unmarshal log entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
