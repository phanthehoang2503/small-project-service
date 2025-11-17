package message

type LogEvent struct {
	Level   string `json:"level"`
	Service string `json:"service"`
	Message string `json:"message"`
	Time    string `json:"time"`
}
