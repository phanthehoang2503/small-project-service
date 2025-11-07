package model

import "time"

type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`          // event time
	Level     string                 `json:"level"`              // debug|info|warn|error
	Service   string                 `json:"service"`            // e.g., auth-service
	Message   string                 `json:"message"`            // message
	TraceID   string                 `json:"trace_id,omitempty"` // trace/request id
	Fields    map[string]interface{} `json:"fields,omitempty"`   // structured fields
}
