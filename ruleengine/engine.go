package ruleengine

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"sentinelx/models"
)

var (
	rules         []models.Rule
	eventCounter  = make(map[string]int)
	lastSeen      = make(map[string]time.Time)
	lastAlertTime = make(map[string]time.Time)
	mutex         sync.RWMutex
)

// ===============================
// LOAD RULES
// ===============================
func LoadRules() error {
	data, err := os.ReadFile("rules/rules.yaml")
	if err != nil {
		return err
	}

	var ruleSet models.RuleSet

	err = yaml.Unmarshal(data, &ruleSet)
	if err != nil {
		return err
	}

	mutex.Lock()
	rules = ruleSet.Rules
	mutex.Unlock()

	fmt.Println("Rule engine loaded rules:", len(ruleSet.Rules))
	return nil
}

// ===============================
// MAIN PROCESSOR
// ===============================
func ProcessEvent(event models.SecurityEvent) *models.Alert {
	mutex.Lock()
	defer mutex.Unlock()

	// ✅ ENSURE METADATA EXISTS
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}

	// 🔥 DETECT ATTACK TYPE
	detectedType := detectAttackType(event)

	// store detected type
	event.Metadata["detected_type"] = detectedType

	// override event type
	event.EventType = detectedType

	// ===============================
	// APPLY RULES
	// ===============================
	for _, rule := range rules {
		if rule.EventType != event.EventType {
			continue
		}

		key := rule.Name + ":" + event.SourceIP

		// ⏱ RESET WINDOW
		if time.Since(lastSeen[key]) > 1*time.Minute {
			eventCounter[key] = 0
		}

		eventCounter[key]++
		lastSeen[key] = time.Now()

		if eventCounter[key] < rule.Threshold {
			continue
		}

		// reset after trigger
		eventCounter[key] = 0

		// 🛑 COOLDOWN (ANTI-SPAM)
		if time.Since(lastAlertTime[key]) < 2*time.Minute {
			return nil
		}

		lastAlertTime[key] = time.Now()

		alert := buildAlert(rule, event)
		return &alert
	}

	return nil
}

// ===============================
// ATTACK DETECTION ENGINE
// ===============================
func detectAttackType(event models.SecurityEvent) string {

	payload := ""

	// ✅ SAFE CAST (CRITICAL FIX)
	if v, ok := event.Metadata["path"].(string); ok {
		payload += v
	}

	if v, ok := event.Metadata["payload"].(string); ok {
		payload += v
	}

	payload = strings.ToLower(payload)

	switch {
	case strings.Contains(payload, "<script>"),
		strings.Contains(payload, "javascript:"),
		strings.Contains(payload, "onerror="):
		return "xss"

	case strings.Contains(payload, "or 1=1"),
		strings.Contains(payload, "' or '1'='1"),
		strings.Contains(payload, "union select"),
		strings.Contains(payload, "drop table"),
		strings.Contains(payload, "select * from"):
		return "sql_injection"

	case strings.Contains(payload, "../"),
		strings.Contains(payload, "..\\"):
		return "path_traversal"

	case strings.Contains(payload, "bash"),
		strings.Contains(payload, "chmod"),
		strings.Contains(payload, "exec"),
		strings.Contains(payload, "curl http"),
		strings.Contains(payload, "wget http"):
		return "rce"

	case strings.Contains(payload, "login"),
		strings.Contains(payload, "password"):
		return "brute_force"

	default:
		return strings.ToLower(event.EventType)
	}
}

// ===============================
// ALERT BUILDER
// ===============================
func buildAlert(rule models.Rule, event models.SecurityEvent) models.Alert {
	score := models.ThreatScoreFromSeverity(rule.Severity)

	return models.Alert{
		ID:          generateRuleAlertID(),
		TenantID:    event.TenantID,
		Timestamp:   time.Now().UTC(),
		Type:        event.EventType, // ✅ normalized type
		Severity:    rule.Severity,
		SourceIP:    event.SourceIP,
		Description: fmt.Sprintf("%s detected from %s", event.EventType, event.SourceIP),
		ThreatScore: score,
		Status:      models.AlertStatusNew,
		Metadata: map[string]interface{}{
			"rule":          rule.Name,
			"event_type":    event.EventType,
			"detected_type": event.Metadata["detected_type"],
		},
	}
}

func generateRuleAlertID() string {
	return "ALT-" + time.Now().UTC().Format("20060102150405.000000000")
}