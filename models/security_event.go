package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SecurityEvent struct {
	EventID     string                 `json:"event_id"`
	TenantID    string                 `json:"tenant_id"`
	Timestamp   int64                  `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	SourceIP    string                 `json:"source_ip"`
	SourcePort  int                    `json:"source_port,omitempty"`
	DestPort    int                    `json:"dest_port,omitempty"`
	Protocol    string                 `json:"protocol"`
	PayloadSize int                    `json:"payload_size,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Alert structure used across detection, storage, websocket, API, and incident handling.
type Alert struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	SourceIP    string                 `json:"source_ip,omitempty"`
	Target      string                 `json:"target,omitempty"`
	Description string                 `json:"description,omitempty"`
	ThreatScore float64                `json:"threat_score"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	TenantID    string                 `json:"tenant_id"`
}

const (
	SeverityLow      = "low"
	SeverityMedium   = "medium"
	SeverityHigh     = "high"
	SeverityCritical = "critical"

	AlertStatusNew       = "new"
	AlertStatusProcessed = "processed"
)

func ThreatScoreFromSeverity(severity string) float64 {
	switch severity {
	case SeverityCritical:
		return 0.95
	case SeverityHigh:
		return 0.75
	case SeverityMedium:
		return 0.55
	default:
		return 0.30
	}
}

func NewSecurityEvent(eventType string) SecurityEvent {
	return SecurityEvent{
		EventID:   uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		EventType: eventType,
		Metadata:  make(map[string]interface{}),
	}
}

func (e *SecurityEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
