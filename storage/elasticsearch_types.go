package storage

// ==============================
// INDEX NAMES
// ==============================

const (
	IndexSecurityEvents  = "security_events"
	IndexAlerts          = "alerts"
	IndexResponseActions = "response_actions"
)

// ==============================
// ELASTICSEARCH CONFIG
// ==============================

type ElasticsearchConfig struct {
	Enabled   bool
	Addresses []string
	Username  string
	Password  string
}

// ==============================
// INDEX MAPPINGS (MULTI-TENANT READY)
// ==============================

// 🔴 ALERTS INDEX
var AlertMapping = `{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "tenant_id": { "type": "keyword" },
      "timestamp": { "type": "date" },
      "event_type": { "type": "keyword" },
      "source_ip": { "type": "ip" },
      "severity": { "type": "keyword" },
      "threat_score": { "type": "float" }
    }
  }
}`

// 🟡 SECURITY EVENTS INDEX
var SecurityEventMapping = `{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "tenant_id": { "type": "keyword" },
      "timestamp": { "type": "date" },
      "event_type": { "type": "keyword" },
      "source_ip": { "type": "ip" },
      "message": { "type": "text" }
    }
  }
}`

// 🔵 RESPONSE ACTIONS INDEX
var ResponseActionMapping = `{
  "mappings": {
    "properties": {
      "id": { "type": "keyword" },
      "tenant_id": { "type": "keyword" },
      "timestamp": { "type": "date" },
      "action_type": { "type": "keyword" },
      "target": { "type": "keyword" },
      "status": { "type": "keyword" }
    }
  }
}`