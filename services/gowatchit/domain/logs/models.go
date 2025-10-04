package logs

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// LogEntry represents a structured log entry for frontend display
type LogEntry struct {
	Time   string                 `json:"time"`
	Level  string                 `json:"level"`
	Source Source                 `json:"source"`
	Msg    string                 `json:"msg"`
	Error  string                 `json:"error,omitempty"`
	Extra  map[string]interface{} `json:"Extra"`
}

// Source represents the source location of a log entry
type Source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle zap's log format
func (l *LogEntry) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to get all fields
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Initialize Extra map with all fields
	l.Extra = make(map[string]interface{})
	for k, v := range rawMap {
		l.Extra[k] = v
	}

	// Extract and convert timestamp (zap uses "ts" field as unix timestamp)
	if ts, ok := rawMap["ts"].(float64); ok {
		l.Time = time.Unix(int64(ts), 0).Format(time.RFC3339)
		delete(l.Extra, "ts")
	} else if tsStr, ok := rawMap["time"].(string); ok {
		l.Time = tsStr
		delete(l.Extra, "time")
	}

	// Extract level (zap uses lowercase, frontend expects uppercase)
	if level, ok := rawMap["level"].(string); ok {
		l.Level = strings.ToUpper(level)
		delete(l.Extra, "level")
	}

	// Extract message
	if msg, ok := rawMap["msg"].(string); ok {
		l.Msg = msg
		delete(l.Extra, "msg")
	}

	// Extract error if present
	if errStr, ok := rawMap["error"].(string); ok {
		l.Error = errStr
		delete(l.Extra, "error")
	}

	// Extract caller/source information (zap format: "file.go:42")
	if caller, ok := rawMap["caller"].(string); ok {
		parts := strings.Split(caller, ":")
		if len(parts) >= 2 {
			l.Source.File = parts[0]
			if line, err := strconv.Atoi(parts[1]); err == nil {
				l.Source.Line = line
			}
		}
		delete(l.Extra, "caller")
	}

	// Extract function name if present
	if fn, ok := rawMap["function"].(string); ok {
		l.Source.Function = fn
		delete(l.Extra, "function")
	}

	// Remove stacktrace from Extra as it's too verbose
	delete(l.Extra, "stacktrace")

	return nil
}
