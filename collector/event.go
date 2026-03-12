package collector

import (
	"encoding/json"   // Used to convert Go structs into JSON format
	"time"            // Used to get the current time for event timestamp

	"github.com/google/uuid" // Used to generate unique event IDs
)

// SecurityEvent represents one security-related event in the system
type SecurityEvent struct {

	// Unique identifier for each event
	EventID string `json:"event_id"`

	// Timestamp when the event happened (in nanoseconds)
	Timestamp int64 `json:"timestamp"`

	// Type of event (example: login_attempt, port_scan, malware_detected)
	EventType string `json:"event_type"`

	// IP address where the event originated
	SourceIP string `json:"source_ip,omitempty"`

	// Destination IP address targeted
	DestIP string `json:"dest_ip,omitempty"`

	// Source port used in the connection
	SourcePort int `json:"source_port,omitempty"`

	// Destination port targeted
	DestPort int `json:"dest_port,omitempty"`

	// Network protocol used (TCP, UDP, HTTP etc.)
	Protocol string `json:"protocol,omitempty"`

	// Size of the network payload in bytes
	PayloadSize int `json:"payload_size,omitempty"`

	// Additional flexible key-value data
	// Example: {"user":"admin", "status":"failed"}
	Metadata map[string]string `json:"metadata,omitempty"`
}

// NewSecurityEvent creates a new security event with default values
func NewSecurityEvent(eventType string) SecurityEvent {

	return SecurityEvent{
		EventID:   uuid.New().String(),   // Generate a unique ID for the event
		Timestamp: time.Now().UnixNano(), // Capture current time in nanoseconds
		EventType: eventType,             // Assign the event type passed by user
		Metadata:  make(map[string]string), // Initialize empty metadata map
	}
}

// ToJSON converts the SecurityEvent struct into JSON format
// This is useful for sending events over network or saving in logs
func (e *SecurityEvent) ToJSON() ([]byte, error) {

	return json.Marshal(e) // Convert struct to JSON byte array
}

// Validate checks if the event contains minimum required data
func (e *SecurityEvent) Validate() bool {

	// EventType must exist
	if e.EventType == "" {
		return false
	}

	// Timestamp must exist
	if e.Timestamp == 0 {
		return false
	}

	// If both checks pass, event is valid
	return true
}