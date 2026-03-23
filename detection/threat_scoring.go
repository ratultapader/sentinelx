package detection

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	WeightAnomalyScore      = 0.4
	WeightIPReputation      = 0.3
	WeightSignatureMatch    = 0.2
	WeightBehaviorDeviation = 0.1
)

type ThreatScorer struct{}

func NewThreatScorer() *ThreatScorer {
	return &ThreatScorer{}
}

func (s *ThreatScorer) Score(signals ThreatSignals) ThreatScoreResult {
	anomaly := clamp01(signals.AnomalyScore)
	ipRep := clamp01(signals.IPReputation)
	signature := clamp01(signals.SignatureMatch)
	behavior := clamp01(signals.BehaviorDeviation)

	score := (anomaly * WeightAnomalyScore) +
		(ipRep * WeightIPReputation) +
		(signature * WeightSignatureMatch) +
		(behavior * WeightBehaviorDeviation)

	score = round(score, 4)
	severity := severityFromScore(score, signals.BaseSeverity)
	reason := buildThreatReason(anomaly, ipRep, signature, behavior, score)

	return ThreatScoreResult{
		AlertID:           signals.AlertID,
		Timestamp:         time.Now().UTC(),
		EventType:         signals.EventType,
		SourceIP:          signals.SourceIP,
		DestinationIP:     signals.DestinationIP,
		Target:            signals.Target,
		ThreatScore:       score,
		Severity:          severity,
		AnomalyScore:      anomaly,
		IPReputation:      ipRep,
		SignatureMatch:    signature,
		BehaviorDeviation: behavior,
		Reason:            reason,
		Description:       signals.Description,
		Metadata:          copyMetadata(signals.Metadata),
	}
}

func (s *ThreatScorer) ScoreWithContext(signals ThreatSignals) ThreatScoreResult {
	result := s.Score(signals)

	boost := contextualBoost(signals.EventType)
	if boost > 0 {
		result.ThreatScore = round(clamp01(result.ThreatScore+boost), 4)
		result.Severity = severityFromScore(result.ThreatScore, signals.BaseSeverity)
		result.Reason = fmt.Sprintf("%s; contextual_boost=%.2f event_type=%s", result.Reason, boost, signals.EventType)
	}

	return result
}

func severityFromScore(score float64, baseSeverity string) string {
	switch {
	case score >= 0.90:
		return "critical"
	case score >= 0.70:
		return "high"
	case score >= 0.50:
		return "medium"
	case score >= 0.25:
		return "low"
	default:
		if strings.TrimSpace(baseSeverity) != "" {
			return strings.ToLower(strings.TrimSpace(baseSeverity))
		}
		return "info"
	}
}

func buildThreatReason(anomaly, ipRep, signature, behavior, score float64) string {
	return fmt.Sprintf(
		"threat score %.4f from anomaly=%.2f ip_reputation=%.2f signature=%.2f behavior=%.2f",
		score, anomaly, ipRep, signature, behavior,
	)
}

func contextualBoost(eventType string) float64 {
	switch strings.ToLower(strings.TrimSpace(eventType)) {
	case "sql_injection":
		return 0.05
	case "reverse_shell":
		return 0.10
	case "container_escape":
		return 0.10
	case "credential_stuffing":
		return 0.05
	case "ransomware_activity":
		return 0.10
	default:
		return 0
	}
}

func DefaultSignalsFromEvent(alertID, eventType, sourceIP, target, description, baseSeverity string) ThreatSignals {
	signature := 0.2
	anomaly := 0.2
	ipRep := 0.2
	behavior := 0.2

	switch strings.ToLower(strings.TrimSpace(eventType)) {
	case "sql_injection":
		signature = 0.95
		anomaly = 0.75
		behavior = 0.70

	case "brute_force":
		signature = 0.70
		anomaly = 0.80
		behavior = 0.85

	case "port_scan":
		signature = 0.65
		anomaly = 0.70
		behavior = 0.75

	case "reverse_shell":
		signature = 0.95
		anomaly = 0.90
		behavior = 0.95

	case "container_escape":
		signature = 0.90
		anomaly = 0.95
		behavior = 0.95

	case "threat_intel_match":
		ipRep = 0.95
		signature = 0.85
		anomaly = 0.60
		behavior = 0.50

	case "repeated_http_requests":
		signature = 0.55
		anomaly = 0.70
		behavior = 0.85

	case "xss_attack":
		signature = 0.85
		anomaly = 0.70
		behavior = 0.60

	case "dir_traversal":
		signature = 0.85
		anomaly = 0.75
		behavior = 0.70

	case "multi_stage_attack":
		signature = 0.95
		anomaly = 0.90
		ipRep = 0.75
		behavior = 0.95

	case "admin_access_watch":
		signature = 0.80
		anomaly = 0.75
		behavior = 0.80
	}

	return ThreatSignals{
		AlertID:           alertID,
		Timestamp:         time.Now().UTC(),
		EventType:         eventType,
		SourceIP:          sourceIP,
		Target:            target,
		AnomalyScore:      anomaly,
		IPReputation:      ipRep,
		SignatureMatch:    signature,
		BehaviorDeviation: behavior,
		BaseSeverity:      baseSeverity,
		Description:       description,
		Metadata:          map[string]interface{}{},
	}
}

func (r ThreatScoreResult) ToAlertMetadata() map[string]interface{} {
	metadata := copyMetadata(r.Metadata)
	metadata["anomaly_score"] = r.AnomalyScore
	metadata["ip_reputation"] = r.IPReputation
	metadata["signature_match"] = r.SignatureMatch
	metadata["behavior_deviation"] = r.BehaviorDeviation
	metadata["score_reason"] = r.Reason
	return metadata
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func round(v float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	return math.Round(v*pow) / pow
}

func copyMetadata(in map[string]interface{}) map[string]interface{} {
	if in == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}