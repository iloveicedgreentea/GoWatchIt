package query

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
)

type GetLogsQuery struct{}

type GetLogsHandler struct{}

func NewGetLogsHandler() (*GetLogsHandler, error) {
	return &GetLogsHandler{}, nil
}

func (h *GetLogsHandler) Handle(ctx context.Context, query *GetLogsQuery) ([]string, error) {
	logFilePath := getLogFilePath()

	file, err := os.Open(logFilePath) // nolint:gosec // We want to read the log file
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Printf("failed to close log file: %v\n", err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	maxLines := 1000
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	var nonEmptyLines []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			nonEmptyLines = append(nonEmptyLines, line)
		}
	}

	return nonEmptyLines, nil
}

func getLogFilePath() string {
	baseDir := os.Getenv("BASE_DIR")
	if baseDir == "" {
		baseDir = "./"
	}
	return baseDir + "/app.log"
}
