package investigation

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Builder struct{}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(sourceIP string, events []TimelineEvent) AttackTimeline {
	fmt.Println("?? TIMELINE BUILDER CALLED")
	filtered := make([]TimelineEvent, 0, len(events))

	for _, ev := range events {

		// if no sourceIP -> take all
		if sourceIP == "" {
			filtered = append(filtered, ev)
			continue
		}

		// normalize values
		evIP := strings.TrimSpace(ev.SourceIP)
		reqIP := strings.TrimSpace(sourceIP)

		// if empty -> include
		if evIP == "" || reqIP == "" {
			filtered = append(filtered, ev)
			continue
		}

		// match
		if evIP == reqIP || strings.Contains(evIP, reqIP) || strings.Contains(reqIP, evIP) {
			filtered = append(filtered, ev)
		}
	}

	// fallback if nothing matched
	if len(filtered) == 0 {
		filtered = events
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})

	timeline := AttackTimeline{
		SourceIP:   sourceIP,
		EventCount: len(filtered),
		Events:     filtered,
	}

	// ✅ ADD THIS LINE
if len(filtered) > 0 {
	timeline.TenantID = filtered[0].TenantID
}

	if len(filtered) == 0 {
		timeline.Conclusion = "no events found for source ip"
		timeline.RiskLevel = "unknown"
		return timeline
	}

	timeline.StartTime = filtered[0].Timestamp
	timeline.EndTime = filtered[len(filtered)-1].Timestamp
	timeline.Stages = uniqueStages(filtered)
	timeline.AttackType = inferAttackType(filtered)
	timeline.RiskLevel = inferRiskLevel(filtered)
	timeline.Conclusion = buildConclusion(filtered, timeline.AttackType, timeline.RiskLevel)
	timeline.Recommendations = buildRecommendations(filtered, timeline.AttackType, timeline.RiskLevel)

	return timeline
}

func uniqueStages(events []TimelineEvent) []string {
	seen := map[string]struct{}{}
	stages := make([]string, 0)

	for _, ev := range events {
		stage := strings.TrimSpace(ev.Stage)
		if stage == "" {
			continue
		}
		if _, ok := seen[stage]; ok {
			continue
		}
		seen[stage] = struct{}{}
		stages = append(stages, stage)
	}
	return stages
}

func inferAttackType(events []TimelineEvent) string {
	counts := map[string]int{}

	for _, ev := range events {
		et := strings.ToLower(strings.TrimSpace(ev.EventType))
		switch et {
		case "port_scan":
			counts["reconnaissance"]++
		case "sql_injection":
			counts["web_attack"]++
		case "xss", "xss_attack":
			counts["web_attack"]++
		case "brute_force", "credential_stuffing":
			counts["credential_attack"]++
		case "reverse_shell", "container_escape", "runtime_compromise":
			counts["post_exploitation"]++
		case "ddos":
			counts["availability_attack"]++
		}
	}

	priority := []string{
		"post_exploitation",
		"web_attack",
		"credential_attack",
		"availability_attack",
		"reconnaissance",
	}

	bestType := "unknown"
	bestCount := 0

	for _, attackType := range priority {
		if counts[attackType] > bestCount {
			bestType = attackType
			bestCount = counts[attackType]
		}
	}

	return bestType
}

func inferRiskLevel(events []TimelineEvent) string {
	maxScore := 0.0
	maxSeverity := ""

	for _, ev := range events {
		if ev.ThreatScore > maxScore {
			maxScore = ev.ThreatScore
		}
		if severityRank(ev.Severity) > severityRank(maxSeverity) {
			maxSeverity = ev.Severity
		}
	}

	switch {
	case maxScore >= 0.9 || strings.EqualFold(maxSeverity, "critical"):
		return "critical"
	case maxScore >= 0.7 || strings.EqualFold(maxSeverity, "high"):
		return "high"
	case maxScore >= 0.5 || strings.EqualFold(maxSeverity, "medium"):
		return "medium"
	case maxScore > 0:
		return "low"
	default:
		return "unknown"
	}
}

func severityRank(s string) int {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

func buildConclusion(events []TimelineEvent, attackType, riskLevel string) string {
	if len(events) == 0 {
		return "no attack activity reconstructed"
	}

	start := events[0].Timestamp.Format(time.RFC3339)
	end := events[len(events)-1].Timestamp.Format(time.RFC3339)

	return fmt.Sprintf(
		"reconstructed %d events from %s to %s; likely attack_type=%s; assessed risk=%s",
		len(events), start, end, attackType, riskLevel,
	)
}

func buildRecommendations(events []TimelineEvent, attackType, riskLevel string) []string {
	recs := []string{
		"review full event trail for the source ip",
		"preserve logs and response execution records",
	}

	switch attackType {
	case "reconnaissance":
		recs = append(recs, "tighten exposure of public services and scan detection thresholds")
	case "web_attack":
		recs = append(recs, "review WAF rules and vulnerable web endpoints")
	case "credential_attack":
		recs = append(recs, "enforce MFA and review authentication protections")
	case "post_exploitation":
		recs = append(recs, "inspect host/container compromise indicators immediately")
	case "availability_attack":
		recs = append(recs, "review upstream rate limiting and DDoS controls")
	}

	if riskLevel == "critical" || riskLevel == "high" {
		recs = append(recs, "escalate to incident response and contain affected assets")
	}

	return recs
}

func NormalizeAlert(doc map[string]interface{}) TimelineEvent {
	event := TimelineEvent{
		ID:          getString(doc, "id"),
		TenantID:    getString(doc, "tenant_id"), // ? ADD HERE
		Timestamp:   getTime(doc, "timestamp"),
		SourceIP:    getString(doc, "source_ip"),
		EventType:   getString(doc, "event_type"),
		Severity:    getString(doc, "severity"),
		ThreatScore: getFloat(doc, "threat_score"),

		// ?? ADD THIS (MITRE)
		MitreTactic:      getString(doc, "mitre_tactic"),
		MitreTechnique:   getString(doc, "mitre_technique"),
		MitreTechniqueID: getString(doc, "mitre_technique_id"),

		Stage:   "detection",
		Source:  "alert",
		Summary: getString(doc, "type"),
		Raw:     doc,
	}

	fmt.Println("DEBUG >>> ALERT:", doc)
	fmt.Println("DEBUG >>> EVENT TYPE:", getString(doc, "event_type"))
	fmt.Println("DEBUG >>> BUILD EVENT:", event)

	return event
}
func NormalizeSecurityEvent(doc map[string]interface{}) TimelineEvent {
	return TimelineEvent{
		ID:          firstNonEmpty(getString(doc, "id"), getString(doc, "event_id")),
		TenantID:    getString(doc, "tenant_id"), // ✅ ADD HERE
		Timestamp:   getTime(doc, "timestamp"),
		SourceIP:    getString(doc, "source_ip"),
		EventType:   getString(doc, "event_type"),
		Severity:    getString(doc, "severity"),
		ThreatScore: getFloat(doc, "threat_score"),
		Stage:       inferStage(getString(doc, "event_type"), "security_event"),
		Source:      "security_event",
		Summary:     firstNonEmpty(getString(doc, "message"), getString(doc, "event_type")),
		Raw:         doc,
	}
}

func NormalizeResponseAction(doc map[string]interface{}) TimelineEvent {
	return TimelineEvent{
		ID:          getString(doc, "id"),
		TenantID:    getString(doc, "tenant_id"), // ✅ ADD HERE
		Timestamp:   getTime(doc, "timestamp"),
		SourceIP:    getString(doc, "source_ip"),
		EventType:   getString(doc, "action_type"),
		Severity:    getString(doc, "severity"),
		ThreatScore: getFloat(doc, "threat_score"),
		Stage:       "response_decision",
		Source:      "response_action",
		Summary:     firstNonEmpty(getString(doc, "reason"), getString(doc, "action_type")),
		Raw:         doc,
	}
}

func NormalizeResponseResult(doc map[string]interface{}) TimelineEvent {
	return TimelineEvent{
		ID:          getString(doc, "id"),
			TenantID:    getString(doc, "tenant_id"), // ✅ ADD HERE
		Timestamp:   getTime(doc, "timestamp"),
		SourceIP:    getString(doc, "source_ip"),
		EventType:   firstNonEmpty(getString(doc, "event"), getString(doc, "action_type")),
		Severity:    getString(doc, "severity"),
		ThreatScore: getFloat(doc, "threat_score"),
		Stage:       "mitigation",
		Source:      "response_result",
		Summary:     firstNonEmpty(getString(doc, "message"), getString(doc, "event")),
		Raw:         doc,
	}
}

func inferStage(eventType, source string) string {
	et := strings.ToLower(strings.TrimSpace(eventType))

	switch et {
	case "port_scan":
		return "reconnaissance"
	case "sql_injection", "xss", "xss_attack", "dir_traversal":
		return "initial_access"
	case "brute_force", "credential_stuffing":
		return "credential_access"
	case "reverse_shell", "container_escape", "runtime_compromise":
		return "execution"
	case "ip_block", "rate_limit", "container_restart", "kubernetes_isolation":
		return "response_decision"
	case "ip_blocked", "rate_limit_applied", "pod_restarted", "kubernetes_isolated":
		return "mitigation"
	default:
		if source == "alert" {
			return "detection"
		}
		return "activity"
	}
}

func getString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func getFloat(m map[string]interface{}, key string) float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}

	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	default:
		return 0
	}
}

func getTime(m map[string]interface{}, key string) time.Time {
	v, ok := m[key]
	if !ok || v == nil {
		return time.Time{}
	}

	switch x := v.(type) {
	case string:
		formats := []string{
			time.RFC3339Nano,
			time.RFC3339,
			"2006-01-02 15:04:05",
		}
		for _, f := range formats {
			if t, err := time.Parse(f, x); err == nil {
				return t
			}
		}
	case time.Time:
		return x
	}

	return time.Time{}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
