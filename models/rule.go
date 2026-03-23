package models

type Rule struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description" json:"description"`
	EventType   string `yaml:"event_type"`
	Threshold   int    `yaml:"threshold"`
	Severity    string `yaml:"severity"`
}

type RuleSet struct {
	Rules []Rule `yaml:"rules"`
}
