package response

import "time"

type KubernetesResult struct {
	ID         string                 `json:"id"`
	ActionID   string                 `json:"action_id"`
	AlertID    string                 `json:"alert_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Event      string                 `json:"event"`
	ActionType string                 `json:"action_type"`
	Namespace  string                 `json:"namespace,omitempty"`
	PodName    string                 `json:"pod_name,omitempty"`
	NodeName   string                 `json:"node_name,omitempty"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}