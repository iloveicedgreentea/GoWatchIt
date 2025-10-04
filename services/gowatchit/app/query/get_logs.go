package query

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/iloveicedgreentea/gowatchit/pkg/logger"
	"github.com/iloveicedgreentea/gowatchit/services/gowatchit/domain/logs"
	"go.uber.org/zap"
)

type GetLogsQuery struct{}

type GetLogsHandler struct{}

func NewGetLogsHandler() (*GetLogsHandler, error) {
	return &GetLogsHandler{}, nil
}

func (h *GetLogsHandler) Handle(ctx context.Context, query *GetLogsQuery) ([]logs.LogEntry, error) {
	log := logger.GetLoggerFromContext(ctx)
	logFilePath := logger.GetLogFilePath()

	file, err := os.Open(logFilePath) // nolint:gosec // We want to read the log file
	if err != nil {
		if os.IsNotExist(err) {
			return []logs.LogEntry{}, nil
		}
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Error("Failed to close log file", zap.Error(closeErr))
		}
	}()

	var entries []logs.LogEntry
	scanner := bufio.NewScanner(file)

	// Create a larger buffer for scanner to handle long lines (like stacktraces)
	const maxCapacity = 512 * 1024 // 512KB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// Read each line and parse as JSON
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var entry logs.LogEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			// Skip malformed lines instead of failing
			log.Debug("Skipping malformed log line", zap.Error(err), zap.ByteString("line", line))
			continue
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %w", err)
	}

	// Return last 1000 entries
	maxEntries := 1000
	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	return entries, nil
}
