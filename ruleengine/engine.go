package ruleengine

import (
	"fmt"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	"sentinelx/models"
)

var (
	rules        []models.Rule
	eventCounter = make(map[string]int)
	mutex        sync.RWMutex
)

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

func GetRules() []models.Rule {
	mutex.RLock()
	defer mutex.RUnlock()

	out := make([]models.Rule, len(rules))
	copy(out, rules)
	return out
}

// ProcessEvent evaluates a security event against loaded YAML rules.
// If a rule matches, it returns a fully built alert.
func ProcessEvent(event models.SecurityEvent) *models.Alert {
	mutex.Lock()
	defer mutex.Unlock()

	for _, rule := range rules {
		if rule.EventType != event.EventType {
			continue
		}

		key := rule.Name + ":" + event.SourceIP
		eventCounter[key]++

		if eventCounter[key] >= rule.Threshold {
			eventCounter[key] = 0

			score := models.ThreatScoreFromSeverity(rule.Severity)

			alert := models.Alert{
				ID:          generateRuleAlertID(),
				Timestamp:   time.Now().UTC(),
				Type:        normalizeRuleName(rule.Name),
				Severity:    rule.Severity,
				SourceIP:    event.SourceIP,
				Description: buildRuleDescription(rule, event),
				ThreatScore: score,
				Status:      models.AlertStatusNew,
				Metadata: map[string]interface{}{
					"rule_name":   rule.Name,
					"event_type":  rule.EventType,
					"threshold":   rule.Threshold,
					"match_count": rule.Threshold,
				},
			}

			return &alert
		}
	}

	return nil
}

func generateRuleAlertID() string {
	return "ALT-" + time.Now().UTC().Format("20060102150405.000000000")
}

func normalizeRuleName(name string) string {
	switch name {
	case "repeated_http_requests":
		return "repeated_http_requests"
	case "admin_access_watch":
		return "admin_access_watch"
	default:
		return name
	}
}

func buildRuleDescription(rule models.Rule, event models.SecurityEvent) string {
	if rule.Description != "" {
		return rule.Description
	}

	return fmt.Sprintf("Rule matched: %s for event type %s from source %s", rule.Name, rule.EventType, event.SourceIP)
}