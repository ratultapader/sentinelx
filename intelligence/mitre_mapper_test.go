package intelligence

import "testing"

func TestMapSQLInjection(t *testing.T) {
	mapper := NewMitreMapper()

	result := mapper.MapEvent("sql_injection")

	if result.Tactic != "Initial Access" {
		t.Fatalf("expected tactic 'Initial Access', got '%s'", result.Tactic)
	}

	if result.TechniqueID != "T1190" {
		t.Fatalf("expected technique ID 'T1190', got '%s'", result.TechniqueID)
	}
}

func TestUnknownEvent(t *testing.T) {
	mapper := NewMitreMapper()

	result := mapper.MapEvent("unknown_attack")

	if result.Tactic != "Unknown" {
		t.Fatalf("expected 'Unknown', got '%s'", result.Tactic)
	}
}

func TestCaseInsensitive(t *testing.T) {
	mapper := NewMitreMapper()

	result := mapper.MapEvent("SQL_INJECTION")

	if result.Tactic != "Initial Access" {
		t.Fatalf("case normalization failed")
	}
}