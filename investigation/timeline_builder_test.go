package investigation

import (
	"testing"
	"time"
)

func TestBuildTimeline(t *testing.T) {
	builder := NewBuilder()

	base := time.Date(2026, 3, 20, 10, 0, 0, 0, time.UTC)
	events := []TimelineEvent{
		{
			ID:          "1",
			Timestamp:   base.Add(2 * time.Minute),
			SourceIP:    "192.168.1.5",
			EventType:   "sql_injection",
			Severity:    "critical",
			ThreatScore: 0.91,
			Stage:       "initial_access",
			Source:      "alert",
			Summary:     "sql injection detected",
		},
		{
			ID:          "2",
			Timestamp:   base,
			SourceIP:    "192.168.1.5",
			EventType:   "port_scan",
			Severity:    "medium",
			ThreatScore: 0.55,
			Stage:       "reconnaissance",
			Source:      "security_event",
			Summary:     "port scan detected",
		},
		{
			ID:          "3",
			Timestamp:   base.Add(7 * time.Minute),
			SourceIP:    "192.168.1.5",
			EventType:   "ip_blocked",
			Severity:    "critical",
			ThreatScore: 0.91,
			Stage:       "mitigation",
			Source:      "response_result",
			Summary:     "ip blocked",
		},
	}

	timeline := builder.Build("192.168.1.5", events)

	if timeline.EventCount != 3 {
		t.Fatalf("expected 3 events, got %d", timeline.EventCount)
	}
	if timeline.Events[0].EventType != "port_scan" {
		t.Fatalf("expected first event to be port_scan, got %s", timeline.Events[0].EventType)
	}
	if timeline.Events[2].EventType != "ip_blocked" {
		t.Fatalf("expected last event to be ip_blocked, got %s", timeline.Events[2].EventType)
	}
	if timeline.RiskLevel != "critical" {
		t.Fatalf("expected risk level critical, got %s", timeline.RiskLevel)
	}
}

func TestNormalizeAlert(t *testing.T) {
	doc := map[string]interface{}{
		"id":           "alert_1",
		"timestamp":    "2026-03-20T10:02:00Z",
		"source_ip":    "192.168.1.5",
		"event_type":   "sql_injection",
		"severity":     "critical",
		"threat_score": 0.95,
		"description":  "SQL injection payload detected",
	}

	ev := NormalizeAlert(doc)

	if ev.Source != "alert" {
		t.Fatalf("expected source alert, got %s", ev.Source)
	}
	if ev.Stage != "initial_access" {
		t.Fatalf("expected stage initial_access, got %s", ev.Stage)
	}
	if ev.EventType != "sql_injection" {
		t.Fatalf("expected sql_injection, got %s", ev.EventType)
	}
}

func TestInferAttackType(t *testing.T) {
	events := []TimelineEvent{
		{EventType: "port_scan"},
		{EventType: "sql_injection"},
		{EventType: "xss"},
	}
	got := inferAttackType(events)
	if got != "web_attack" {
		t.Fatalf("expected web_attack, got %s", got)
	}
}