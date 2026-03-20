package ruleengine

import (
	"fmt"
	"os"
	"sync"

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

func ProcessEvent(event models.SecurityEvent) *models.Rule {
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

			matchedRule := rule
			return &matchedRule
		}
	}

	return nil
}