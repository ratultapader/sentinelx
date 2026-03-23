package storage

import "testing"

func TestNeo4jConstants(t *testing.T) {
	if NodeAttackerIP == "" {
		t.Fatal("NodeAttackerIP must not be empty")
	}
	if RelAttacked == "" {
		t.Fatal("RelAttacked must not be empty")
	}
}

func TestAttackGraphRecordValidationShape(t *testing.T) {
	record := AttackGraphRecord{
		AlertID:     "alert_1",
		SourceIP:    "192.168.1.10",
		Server:      "payments-api",
		APIEndpoint: "/login",
		EventType:   "sql_injection",
		Severity:    "critical",
		ThreatScore: 0.95,
	}
	if record.AlertID == "" {
		t.Fatal("expected alert id")
	}
	if record.SourceIP == "" {
		t.Fatal("expected source ip")
	}
}