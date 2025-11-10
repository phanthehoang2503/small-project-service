package message

type LogEvent struct {
	Service string      `json:"service"` // which service sent this
	Level   string      `json:"level"`   // info, warn, error
	Message string      `json:"message"` // description
	Meta    interface{} `json:"meta,omitempty"`
	Time    string      `json:"time,omitempty"`
}
