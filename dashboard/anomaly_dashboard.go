package dashboard

import (
	"context"
	"time"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// ✅ FINAL RESPONSE STRUCTURE (WITH MITRE)
type DashboardResponse struct {
	Summary struct {
		TotalAlerts     int     `json:"total_alerts"`
		CriticalAlerts  int     `json:"critical_alerts"`
		HighAlerts      int     `json:"high_alerts"`
		MediumAlerts    int     `json:"medium_alerts"`
		LowAlerts       int     `json:"low_alerts"`
		UniqueAttackers int     `json:"unique_attackers"`
		AvgThreatScore  float64 `json:"avg_threat_score"`
	} `json:"summary"`

	AttackFrequency []struct {
		EventType string `json:"event_type"`
		Count     int    `json:"count"`
	} `json:"attack_frequency"`

	ThreatDistribution []struct {
		Bucket string `json:"bucket"`
		Count  int    `json:"count"`
	} `json:"threat_distribution"`

	TopEndpoints []struct {
		Target string `json:"target"`
		Count  int    `json:"count"`
	} `json:"top_attacked_endpoints"`

	AnomalyTrend []struct {
		Time  string  `json:"time"`
		Score float64 `json:"avg_anomaly_score"`
	} `json:"anomaly_trend"`

	// 🔥 NEW MITRE FIELD
	MitreDistribution []struct {
		Tactic string `json:"tactic"`
		Count  int    `json:"count"`
	} `json:"mitre_distribution"`
}

// 🔥 CORE BUILDER
func (s *DashboardService) BuildDashboard(ctx context.Context, alerts []map[string]interface{}) (*DashboardResponse, error) {

	resp := &DashboardResponse{}

	uniqueIPs := make(map[string]bool)
	attackTypeCount := make(map[string]int)
	endpointMap := make(map[string]int)
	mitreCounts := make(map[string]int) // 🔥 NEW

	// 🔥 Trend + Threat
	trend := map[string][]float64{}
	buckets := map[string]int{
		"0-0.2":   0,
		"0.2-0.4": 0,
		"0.4-0.6": 0,
		"0.6-0.8": 0,
		"0.8-1":   0,
	}

	totalThreat := 0.0
	threatCount := 0

	for _, alert := range alerts {

		resp.Summary.TotalAlerts++

		// 🔹 Severity
		if sev, ok := alert["severity"].(string); ok {
			switch sev {
			case "critical":
				resp.Summary.CriticalAlerts++
			case "high":
				resp.Summary.HighAlerts++
			case "medium":
				resp.Summary.MediumAlerts++
			case "low":
				resp.Summary.LowAlerts++
			}
		}

		// 🔹 Unique attackers
		if ip, ok := alert["source_ip"].(string); ok {
			uniqueIPs[ip] = true
		}

		// 🔹 Attack types
		if t, ok := alert["type"].(string); ok && t != "" {
			attackTypeCount[t]++
		}

		// 🔥 MITRE aggregation (FIXED)
meta := getMap(alert, "metadata")

if tactic, ok := meta["mitre_tactic"].(string); ok && tactic != "" {
	mitreCounts[tactic]++
}

		// 🔹 Endpoints
		target := getString(alert, "target")
		if target != "" {
			endpointMap[target]++
		}

		// 🔹 Threat score
		score := getFloat(alert, "threat_score")
		if score > 0 {
			totalThreat += score
			threatCount++
		}

		// 🔹 Buckets
		switch {
		case score < 0.2:
			buckets["0-0.2"]++
		case score < 0.4:
			buckets["0.2-0.4"]++
		case score < 0.6:
			buckets["0.4-0.6"]++
		case score < 0.8:
			buckets["0.6-0.8"]++
		default:
			buckets["0.8-1"]++
		}

		// 🔹 Anomaly trend
		ts := getString(alert, "timestamp")
		meta = getMap(alert, "metadata")
		anomaly := getFloat(meta, "anomaly_score")

		if ts != "" {
			t, err := time.Parse(time.RFC3339, ts)
			if err == nil {
				key := t.Format("15:00")
				trend[key] = append(trend[key], anomaly)
			}
		}
	}

	// 🔹 Final calculations
	resp.Summary.UniqueAttackers = len(uniqueIPs)

	if threatCount > 0 {
		resp.Summary.AvgThreatScore = totalThreat / float64(threatCount)
	}

	// 🔹 Attack frequency
	for k, v := range attackTypeCount {
		resp.AttackFrequency = append(resp.AttackFrequency, struct {
			EventType string `json:"event_type"`
			Count     int    `json:"count"`
		}{
			EventType: k,
			Count:     v,
		})
	}

	// 🔹 Threat distribution
	for k, v := range buckets {
		resp.ThreatDistribution = append(resp.ThreatDistribution, struct {
			Bucket string `json:"bucket"`
			Count  int    `json:"count"`
		}{
			Bucket: k,
			Count:  v,
		})
	}

	// 🔹 Endpoints
	for k, v := range endpointMap {
		resp.TopEndpoints = append(resp.TopEndpoints, struct {
			Target string `json:"target"`
			Count  int    `json:"count"`
		}{
			Target: k,
			Count:  v,
		})
	}

	// 🔹 Trend
	for k, values := range trend {
		sum := 0.0
		for _, v := range values {
			sum += v
		}

		resp.AnomalyTrend = append(resp.AnomalyTrend, struct {
			Time  string  `json:"time"`
			Score float64 `json:"avg_anomaly_score"`
		}{
			Time:  k,
			Score: sum / float64(len(values)),
		})
	}

	// 🔥 MITRE distribution (FINAL)
	for tactic, count := range mitreCounts {

	if tactic == "" || tactic == "Unknown" {
		continue // 🔥 remove useless data
	}

	resp.MitreDistribution = append(resp.MitreDistribution, struct {
		Tactic string `json:"tactic"`
		Count  int    `json:"count"`
	}{
		Tactic: tactic,
		Count:  count,
	})
}

	return resp, nil
}

//
// 🔥 SAFE HELPERS
//

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

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key].(map[string]interface{}); ok {
		return v
	}
	return map[string]interface{}{}
}