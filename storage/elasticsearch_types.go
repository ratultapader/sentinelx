package storage

const (
	IndexSecurityEvents  = "security_events"
	IndexAlerts          = "alerts"
	IndexResponseActions = "response_actions"
)

type ElasticsearchConfig struct {
	Enabled   bool
	Addresses []string
	Username  string
	Password  string
}