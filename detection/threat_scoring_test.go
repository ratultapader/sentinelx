package detection

import "testing"

func TestThreatScoreFormula(t *testing.T) {
	scorer := NewThreatScorer()

	result := scorer.Score(ThreatSignals{
		AlertID:           "alert_1",
		EventType:         "sql_injection",
		SourceIP:          "192.168.1.10",
		AnomalyScore:      0.9,
		IPReputation:      0.8,
		SignatureMatch:    0.7,
		BehaviorDeviation: 0.6,
	})

	expected := 0.80
	if result.ThreatScore != expected {
		t.Fatalf("expected %.2f, got %.4f", expected, result.ThreatScore)
	}
	if result.Severity != "high" {
		t.Fatalf("expected severity high, got %s", result.Severity)
	}
}

func TestThreatScoreClamp(t *testing.T) {
	scorer := NewThreatScorer()

	result := scorer.Score(ThreatSignals{
		AlertID:           "alert_2",
		EventType:         "weird_event",
		SourceIP:          "192.168.1.20",
		AnomalyScore:      1.5,
		IPReputation:      -1.0,
		SignatureMatch:    2.0,
		BehaviorDeviation: 0.5,
	})

	expected := 0.65
	if result.ThreatScore != expected {
		t.Fatalf("expected %.2f, got %.4f", expected, result.ThreatScore)
	}
	if result.Severity != "medium" {
		t.Fatalf("expected severity medium, got %s", result.Severity)
	}
}

func TestSeverityMapping(t *testing.T) {
	tests := []struct {
		score    float64
		expected string
	}{
		{0.95, "critical"},
		{0.75, "high"},
		{0.55, "medium"},
		{0.30, "low"},
		{0.10, "info"},
	}

	for _, tt := range tests {
		got := severityFromScore(tt.score, "")
		if got != tt.expected {
			t.Fatalf("score %.2f expected %s got %s", tt.score, tt.expected, got)
		}
	}
}

func TestContextualBoost(t *testing.T) {
	scorer := NewThreatScorer()

	result := scorer.ScoreWithContext(ThreatSignals{
		AlertID:           "alert_3",
		EventType:         "reverse_shell",
		SourceIP:          "192.168.1.30",
		AnomalyScore:      0.80,
		IPReputation:      0.80,
		SignatureMatch:    0.80,
		BehaviorDeviation: 0.80,
	})

	if result.ThreatScore != 0.9 {
		t.Fatalf("expected 0.9 got %.4f", result.ThreatScore)
	}
	if result.Severity != "critical" {
		t.Fatalf("expected critical got %s", result.Severity)
	}
}