package storage

import (
	"encoding/json"
	"log"

	"sentinelx/models"
)

// SaveAlert stores a full alert record in the database.
func SaveAlert(alert models.Alert) {
	metadataJSON, err := json.Marshal(alert.Metadata)
	if err != nil {
		log.Println("failed to marshal alert metadata:", err)
		metadataJSON = []byte("{}")
	}

	query := `
	INSERT INTO alerts (
		id,
		timestamp,
		type,
		severity,
		source_ip,
		target,
		description,
		threat_score,
		status,
		metadata
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = DB.Exec(
		query,
		alert.ID,
		alert.Timestamp,
		alert.Type,
		alert.Severity,
		alert.SourceIP,
		alert.Target,
		alert.Description,
		alert.ThreatScore,
		alert.Status,
		string(metadataJSON),
	)
	if err != nil {
		log.Println("failed to save alert:", err)
		return
	}

	IndexAlertDoc(map[string]interface{}{
		"id":           alert.ID,
		"timestamp":    alert.Timestamp,
		"type":         alert.Type,
		"severity":     alert.Severity,
		"source_ip":    alert.SourceIP,
		"target":       alert.Target,
		"description":  alert.Description,
		"threat_score": alert.ThreatScore,
		"status":       alert.Status,
		"metadata":     alert.Metadata,
	}, alert.ID)
}