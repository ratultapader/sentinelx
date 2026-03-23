package storage

import "time"

const (
	NodeAttackerIP     = "AttackerIP"
	NodeServer         = "Server"
	NodeContainer      = "Container"
	NodeAPIEndpoint    = "APIEndpoint"
	NodeAlert          = "Alert"
	NodeResponseAction = "ResponseAction"
)

const (
	RelAttacked    = "ATTACKED"
	RelConnectedTo = "CONNECTED_TO"
	RelExploited   = "EXPLOITED"
	RelTriggered   = "TRIGGERED"
	RelTargeted    = "TARGETED"
	RelMitigatedBy = "MITIGATED_BY"
)

type GraphNode struct {
	Label      string                 `json:"label"`
	Key        string                 `json:"key"`
	Properties map[string]interface{} `json:"properties"`
}

type GraphRelationship struct {
	Type       string                 `json:"type"`
	FromLabel  string                 `json:"from_label"`
	FromKey    string                 `json:"from_key"`
	ToLabel    string                 `json:"to_label"`
	ToKey      string                 `json:"to_key"`
	Properties map[string]interface{} `json:"properties"`
}

type AttackGraphRecord struct {
	AlertID        string    `json:"alert_id"`
	Timestamp      time.Time `json:"timestamp"`
	SourceIP       string    `json:"source_ip"`
	DestinationIP  string    `json:"destination_ip,omitempty"`
	Server         string    `json:"server,omitempty"`
	Container      string    `json:"container,omitempty"`
	APIEndpoint    string    `json:"api_endpoint,omitempty"`
	EventType      string    `json:"event_type"`
	Severity       string    `json:"severity"`
	ThreatScore    float64   `json:"threat_score"`
	ResponseAction string    `json:"response_action,omitempty"`
}

type GraphNodeView struct {
	Label      string                 `json:"label"`
	Key        string                 `json:"key"`
	Properties map[string]interface{} `json:"properties"`
}

type GraphLinkView struct {
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

type GraphView struct {
	Nodes []GraphNodeView `json:"nodes"`
	Links []GraphLinkView `json:"links"`
}