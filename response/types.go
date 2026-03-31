package response

import "time"

const (
	ActionIPBlock          = "ip_block"
	ActionRateLimit        = "rate_limit"
	ActionContainerRestart = "container_restart"
	ActionK8sIsolation     = "kubernetes_isolation"
	ActionAlertOnly        = "alert_only"

	StatusPending  = "pending"
	StatusExecuted = "executed"
	StatusFailed   = "failed"
)

type Action struct {
	ID          string                 `json:"id"`
	AlertID     string                 `json:"alert_id"`
	Timestamp   time.Time              `json:"timestamp"`
	ActionType  string                 `json:"action_type"`
	SourceIP    string                 `json:"source_ip,omitempty"`
	Target      string                 `json:"target,omitempty"`
	TenantID string `json:"tenant_id"` // ✅ ADD
	Severity    string                 `json:"severity"`
	ThreatScore float64                `json:"threat_score"`
	Reason      string                 `json:"reason"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}