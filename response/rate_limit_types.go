package response

import "time"

type RateLimitResult struct {
	ID          string                 `json:"id"`
	ActionID    string                 `json:"action_id"`
	AlertID     string                 `json:"alert_id"`
	Timestamp   time.Time              `json:"timestamp"`
	Event       string                 `json:"event"`
	SourceIP    string                 `json:"source_ip"`
	ActionType  string                 `json:"action_type"`
	Status      string                 `json:"status"`
	Message     string                 `json:"message"`
	LimitPerSec int                    `json:"limit_per_sec"`
	Burst       int                    `json:"burst"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}