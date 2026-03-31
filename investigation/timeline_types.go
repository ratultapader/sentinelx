package investigation

import "time"

// ===============================
// Timeline Event (Enhanced)
// ===============================

type TimelineEvent struct {
	ID          string                 `json:"id"`
	TenantID    string                 `json:"tenant_id,omitempty"` // ✅ ADDED
	Timestamp   time.Time              `json:"timestamp"`
	SourceIP    string                 `json:"source_ip"`
	EventType   string                 `json:"event_type"`
	Severity    string                 `json:"severity,omitempty"`
	ThreatScore float64                `json:"threat_score,omitempty"`

	// =========================
	// 🔥 MITRE INTELLIGENCE
	// =========================
	MitreTactic      string `json:"mitre_tactic,omitempty"`
	MitreTechnique   string `json:"mitre_technique,omitempty"`
	MitreTechniqueID string `json:"mitre_technique_id,omitempty"`

	// =========================
	// Investigation Fields
	// =========================
	Stage   string `json:"stage"`
	Source  string `json:"source"`
	Summary string `json:"summary"`

	Raw map[string]interface{} `json:"raw,omitempty"`
}

// ===============================
// Attack Timeline
// ===============================

type AttackTimeline struct {
	TenantID string `json:"tenant_id,omitempty"` // ✅ ADDED

	SourceIP   string          `json:"source_ip"`
	StartTime  time.Time       `json:"start_time"`
	EndTime    time.Time       `json:"end_time"`
	EventCount int             `json:"event_count"`
	Stages     []string        `json:"stages"`
	Events     []TimelineEvent `json:"events"`

	// Analysis Output
	AttackType      string   `json:"attack_type"`
	RiskLevel       string   `json:"risk_level"`
	Conclusion      string   `json:"conclusion"`
	Recommendations []string `json:"recommendations"`
}