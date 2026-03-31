package reporting

import (
	"encoding/json"
	"strings"
	"time"

	"sentinelx/investigation"
)

type ReportBuilder struct{}

func NewReportBuilder() *ReportBuilder {
	return &ReportBuilder{}
}

func (b *ReportBuilder) Build(
	incidentID string,
	alert map[string]interface{},
	timeline investigation.AttackTimeline,
	actions []map[string]interface{},
	results []map[string]interface{},
) IncidentReport {

	mitreTactic := getString(alert, "mitre_tactic")
	if mitreTactic == "" {
		mitreTactic = "Initial Access"
	}

	mitreTechnique := getString(alert, "mitre_technique")
	if mitreTechnique == "" {
		mitreTechnique = "Suspicious Network Activity"
	}

	t := getString(alert, "type")
	if t == "" {
		t = "Security Incident"
	}

	return IncidentReport{
		IncidentID:  incidentID,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		SourceIP:    getString(alert, "source_ip"),

		Title: "Incident Report: " + t,

		ExecutiveSummary: "Automated incident detected by SentinelX based on abnormal activity patterns.",

		AttackType:  timeline.AttackType,
		RiskLevel:   timeline.RiskLevel,
		Severity:    getString(alert, "severity"),
		ThreatScore: getFloat(alert, "threat_score"),

		MitreTactic:      mitreTactic,
		MitreTechnique:   mitreTechnique,
		MitreTechniqueID: getString(alert, "mitre_technique_id"),

		AttackChain: buildAttackChain(timeline),
		ActionsTaken: buildActions(actions, results),

		RecommendedRemediation: []string{
			"Block the malicious IP",
			"Patch vulnerable endpoints",
			"Enable WAF protection",
			"Review application logs",
		},

		Evidence: []EvidenceItem{
			{
				Type:    "alert",
				Source:  "elasticsearch",
				Summary: getString(alert, "type"),
			},
		},
	}
}

// ================= HELPERS =================

func buildAttackChain(t investigation.AttackTimeline) []AttackChainStep {
	out := []AttackChainStep{}

	for _, ev := range t.Events {
		out = append(out, AttackChainStep{
			Timestamp: ev.Timestamp.Format("2006-01-02 15:04:05"),

			Stage:     safe(ev.Stage, "detection"),
			EventType: safe(ev.EventType, ev.Stage),
			Summary:   safe(ev.Summary, "Suspicious activity detected"),

			Severity:    ev.Severity,
			ThreatScore: ev.ThreatScore,

			MitreTactic:      ev.MitreTactic,
			MitreTechnique:   ev.MitreTechnique,
			MitreTechniqueID: ev.MitreTechniqueID,
		})
	}

	return out
}

func safe(v string, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}

func buildActions(a, r []map[string]interface{}) []ActionTaken {
	out := []ActionTaken{}

	for _, x := range a {
		out = append(out, ActionTaken{
			Timestamp: getString(x, "timestamp"),
			Action:    getString(x, "action_type"),
			Status:    getString(x, "status"),
		})
	}

	return out
}

func (r IncidentReport) ToJSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}

