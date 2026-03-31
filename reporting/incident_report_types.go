package reporting

type IncidentReport struct {
	IncidentID       string  `json:"incident_id"`
	GeneratedAt      string  `json:"generated_at"`
	SourceIP         string  `json:"source_ip"`
	Title            string  `json:"title"`
	ExecutiveSummary string  `json:"executive_summary"`

	AttackType  string  `json:"attack_type"`
	RiskLevel   string  `json:"risk_level"`
	Severity    string  `json:"severity"`
	ThreatScore float64 `json:"threat_score"`

	MitreTactic      string `json:"mitre_tactic,omitempty"`
	MitreTechnique   string `json:"mitre_technique,omitempty"`
	MitreTechniqueID string `json:"mitre_technique_id,omitempty"`

	AttackChain []AttackChainStep `json:"attack_chain"`
	ActionsTaken []ActionTaken    `json:"actions_taken"`

	RecommendedRemediation []string       `json:"recommended_remediation"`
	Evidence               []EvidenceItem `json:"evidence"`
}

type AttackChainStep struct {
	Timestamp   string  `json:"timestamp"`
	Stage       string  `json:"stage"`
	EventType   string  `json:"event_type"`
	Summary     string  `json:"summary"`
	Severity    string  `json:"severity,omitempty"`
	ThreatScore float64 `json:"threat_score,omitempty"`

	MitreTactic      string `json:"mitre_tactic,omitempty"`
	MitreTechnique   string `json:"mitre_technique,omitempty"`
	MitreTechniqueID string `json:"mitre_technique_id,omitempty"`
}

type ActionTaken struct {
	Timestamp string `json:"timestamp"`
	Action    string `json:"action"`
	Status    string `json:"status,omitempty"`
	Reason    string `json:"reason,omitempty"`
}

type EvidenceItem struct {
	Type      string `json:"type"`
	Source    string `json:"source"`
	Reference string `json:"reference"`
	Summary   string `json:"summary"`
}