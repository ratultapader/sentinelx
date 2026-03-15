package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SecurityEvent struct {
	EventID     string            `json:"event_id"`
	Timestamp   int64             `json:"timestamp"`
	EventType   string            `json:"event_type"`
	SourceIP    string            `json:"source_ip"`
	SourcePort  int               `json:"source_port,omitempty"`
	DestPort    int               `json:"dest_port,omitempty"`
	Protocol    string            `json:"protocol"`
	PayloadSize int               `json:"payload_size,omitempty"`
	Metadata    map[string]string `json:"metadata"`
}

// Alert structure moved here to avoid import cycle
type Alert struct {
	Timestamp   int64  `json:"timestamp"`
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	SourceIP    string `json:"source_ip"`
	Description string `json:"description"`
}

func NewSecurityEvent(eventType string) SecurityEvent {

	return SecurityEvent{
		EventID:   uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		EventType: eventType,
		Metadata:  make(map[string]string),
	}
}

func (e *SecurityEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
