package api

type IncidentListItem struct {
	ID          string  `json:"id"`
	Timestamp   string  `json:"timestamp"`
	EventType   string  `json:"event_type"`
	SourceIP    string  `json:"source_ip"`
	Severity    string  `json:"severity"`
	ThreatScore float64 `json:"threat_score"`

	MitreTactic      string `json:"mitre_tactic,omitempty"`
	MitreTechnique   string `json:"mitre_technique,omitempty"`
	MitreTechniqueID string `json:"mitre_technique_id,omitempty"`
}

type IncidentListResponse struct {
	Items []IncidentListItem `json:"items"`
	Count int                `json:"count"`
}

