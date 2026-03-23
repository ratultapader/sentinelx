package detection

import "time"

type ThreatSignals struct {
	AlertID           string                 `json:"alert_id"`
	Timestamp         time.Time              `json:"timestamp"`
	EventType         string                 `json:"event_type"`
	SourceIP          string                 `json:"source_ip"`
	DestinationIP     string                 `json:"destination_ip,omitempty"`
	Target            string                 `json:"target,omitempty"`
	AnomalyScore      float64                `json:"anomaly_score"`
	IPReputation      float64                `json:"ip_reputation"`
	SignatureMatch    float64                `json:"signature_match"`
	BehaviorDeviation float64                `json:"behavior_deviation"`
	BaseSeverity      string                 `json:"base_severity,omitempty"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type ThreatScoreResult struct {
	AlertID           string                 `json:"alert_id"`
	Timestamp         time.Time              `json:"timestamp"`
	EventType         string                 `json:"event_type"`
	SourceIP          string                 `json:"source_ip"`
	DestinationIP     string                 `json:"destination_ip,omitempty"`
	Target            string                 `json:"target,omitempty"`
	ThreatScore       float64                `json:"threat_score"`
	Severity          string                 `json:"severity"`
	AnomalyScore      float64                `json:"anomaly_score"`
	IPReputation      float64                `json:"ip_reputation"`
	SignatureMatch    float64                `json:"signature_match"`
	BehaviorDeviation float64                `json:"behavior_deviation"`
	Reason            string                 `json:"reason"`
	Description       string                 `json:"description,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}