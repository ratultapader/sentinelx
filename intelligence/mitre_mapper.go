package intelligence

import "strings"

// ===============================
// Mapper
// ===============================

type MitreMapper struct {
	mappings map[string]MitreMapping
}

func NewMitreMapper() *MitreMapper {
	return &MitreMapper{
		mappings: defaultMitreMappings(),
	}
}

// ===============================
// Core Mapping Logic (FINAL)
// ===============================

func (m *MitreMapper) MapEvent(eventType string) EnrichedMitreContext {
	normalized := normalizeEventType(eventType)

	// ✅ Direct mapping
	if mapping, ok := m.mappings[normalized]; ok {
		copyMapping := mapping

		return EnrichedMitreContext{
			EventType:    normalized,
			Tactic:       mapping.Tactic,
			Technique:    mapping.Technique,
			TechniqueID:  mapping.TechniqueID,
			SubTechnique: mapping.SubTechnique,
			Confidence:   mapping.Confidence,
			Raw:          &copyMapping,
		}
	}

	// 🔥 AUTO FALLBACK
	switch {

	case strings.Contains(normalized, "scan"):
		return EnrichedMitreContext{
			EventType:   normalized,
			Tactic:      "Reconnaissance",
			Technique:   "Active Scanning",
			TechniqueID: "T1595",
			Confidence:  0.6,
		}

	case strings.Contains(normalized, "request"):
		return EnrichedMitreContext{
			EventType:   normalized,
			Tactic:      "Reconnaissance",
			Technique:   "Gather Victim Network Information",
			TechniqueID: "T1590",
			Confidence:  0.5,
		}

	case strings.Contains(normalized, "brute"):
		return EnrichedMitreContext{
			EventType:   normalized,
			Tactic:      "Credential Access",
			Technique:   "Brute Force",
			TechniqueID: "T1110",
			Confidence:  0.7,
		}

	case strings.Contains(normalized, "inject"):
		return EnrichedMitreContext{
			EventType:   normalized,
			Tactic:      "Initial Access",
			Technique:   "Exploit Public-Facing Application",
			TechniqueID: "T1190",
			Confidence:  0.7,
		}
	}

	// fallback
	return EnrichedMitreContext{
		EventType:   normalized,
		Tactic:      "Unknown",
		Technique:   "Unknown",
		Confidence:  0.0,
	}
}

// ===============================
// Helpers
// ===============================

func normalizeEventType(eventType string) string {
	return strings.ToLower(strings.TrimSpace(eventType))
}

// ===============================
// Default MITRE Mappings
// ===============================

func defaultMitreMappings() map[string]MitreMapping {
	return map[string]MitreMapping{

		"port_scan": {
			EventType:   "port_scan",
			Tactic:      "Reconnaissance",
			Technique:   "Active Scanning",
			TechniqueID: "T1595",
			Confidence:  0.95,
		},

		"repeated_http_requests": {
			EventType:   "repeated_http_requests",
			Tactic:      "Reconnaissance",
			Technique:   "Active Scanning",
			TechniqueID: "T1595",
			Confidence:  0.75,
		},

		"sql_injection": {
			EventType:   "sql_injection",
			Tactic:      "Initial Access",
			Technique:   "Exploit Public-Facing Application",
			TechniqueID: "T1190",
			Confidence:  0.95,
		},

		"xss_attack": {
			EventType:   "xss_attack",
			Tactic:      "Initial Access",
			Technique:   "Exploit Public-Facing Application",
			TechniqueID: "T1190",
			Confidence:  0.80,
		},

		"dir_traversal": {
			EventType:   "dir_traversal",
			Tactic:      "Initial Access",
			Technique:   "Exploit Public-Facing Application",
			TechniqueID: "T1190",
			Confidence:  0.85,
		},

		"threat_intel_match": {
			EventType:   "threat_intel_match",
			Tactic:      "Reconnaissance",
			Technique:   "Gather Victim Network Information",
			TechniqueID: "T1590",
			Confidence:  0.85,
		},

		"http_request": {
	EventType:   "http_request",
	Tactic:      "Reconnaissance",
	Technique:   "Gather Victim Network Information",
	TechniqueID: "T1590",
	Confidence:  0.6,
},

		"brute_force": {
			EventType:   "brute_force",
			Tactic:      "Credential Access",
			Technique:   "Brute Force",
			TechniqueID: "T1110",
			Confidence:  0.97,
		},
	}
}

// ===============================
// Alert Enrichment
// ===============================

func EnrichAlertMetadata(
	alert map[string]interface{},
	ctx EnrichedMitreContext,
) map[string]interface{} {

	if alert == nil {
		alert = map[string]interface{}{}
	}

	metadata, _ := alert["metadata"].(map[string]interface{})
	if metadata == nil {
		metadata = map[string]interface{}{}
	}

	metadata["mitre_tactic"] = ctx.Tactic
	metadata["mitre_technique"] = ctx.Technique
	metadata["mitre_technique_id"] = ctx.TechniqueID
	metadata["mitre_sub_technique"] = ctx.SubTechnique
	metadata["mitre_confidence"] = ctx.Confidence

	alert["metadata"] = metadata

	// flatten
	alert["mitre_tactic"] = ctx.Tactic
	alert["mitre_technique"] = ctx.Technique
	alert["mitre_technique_id"] = ctx.TechniqueID

	return alert
}