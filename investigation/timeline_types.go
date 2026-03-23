package investigation

import "time"

type TimelineEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	SourceIP    string                 `json:"source_ip"`
	EventType   string                 `json:"event_type"`
	Severity    string                 `json:"severity,omitempty"`
	ThreatScore float64                `json:"threat_score,omitempty"`
	Stage       string                 `json:"stage"`
	Source      string                 `json:"source"`
	Summary     string                 `json:"summary"`
	Raw         map[string]interface{} `json:"raw,omitempty"`
}

type AttackTimeline struct {
	SourceIP        string          `json:"source_ip"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	EventCount      int             `json:"event_count"`
	Stages          []string        `json:"stages"`
	Events          []TimelineEvent `json:"events"`
	AttackType      string          `json:"attack_type"`
	RiskLevel       string          `json:"risk_level"`
	Conclusion      string          `json:"conclusion"`
	Recommendations []string        `json:"recommendations"`
}