package models

import (
	"encoding/json"
)

// models for sending logs to frontend

type LogEntry struct {
	Time   string                 `json:"time"`
	Level  string                 `json:"level"`
	Source Source                 `json:"source"`
	Msg    string                 `json:"msg"`
	Error  string                 `json:"error"`
	Extra  map[string]interface{} `json:",inline"`
}
type Source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func (l *LogEntry) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map to get all fields
	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// Store ALL fields in Extra first
	l.Extra = make(map[string]interface{})
	for k, v := range rawMap {
		l.Extra[k] = v
	}

	// Extract known fields
	l.Time, _ = rawMap["time"].(string)
	l.Level, _ = rawMap["level"].(string)
	l.Msg, _ = rawMap["msg"].(string)
	l.Error, _ = rawMap["error"].(string)

	// Handle source separately as it's a nested structure
	if sourceRaw, ok := rawMap["source"].(map[string]interface{}); ok {
		l.Source = Source{
			Function: sourceRaw["function"].(string),
			File:     sourceRaw["file"].(string),
			Line:     int(sourceRaw["line"].(float64)),
		}
	}

	// Remove known fields from Extra
	delete(l.Extra, "time")
	delete(l.Extra, "level")
	delete(l.Extra, "msg")
	delete(l.Extra, "source")
	delete(l.Extra, "error")

	return nil
}
