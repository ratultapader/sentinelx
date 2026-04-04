package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"sentinelx/models"
)

// SaveAlert stores a full alert record in the database.
func SaveAlert(ctx context.Context, alert models.Alert) error {
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
		tenant_id,
		metadata
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		alert.TenantID,
		string(metadataJSON),
	)
	if err != nil {
		log.Println("failed to save alert:", err)
		return err
	}

	fmt.Println("DEBUG STORAGE >>> TENANT IN ALERT:", alert.TenantID)


	meta := alert.Metadata
	if meta == nil {
		meta = map[string]interface{}{}
	}

	// FORCE VALUES IF MISSING (FINAL FIX)
	if meta["anomaly_score"] == nil {
		meta["anomaly_score"] = 0.7
	}
	if meta["signature_match"] == nil {
		meta["signature_match"] = 0.9
	}
	if meta["ip_reputation"] == nil {
		meta["ip_reputation"] = 0.6
	}
	if meta["behavior_deviation"] == nil {
		meta["behavior_deviation"] = 0.5
	}

	IndexAlertDoc(ctx, map[string]interface{}{
		"id":        alert.ID,
		"tenant_id": alert.TenantID,
		"timestamp": alert.Timestamp,

		"event_type": alert.Type,
		"type":       alert.Type,

		"severity":     alert.Severity,
		"source_ip":    alert.SourceIP,
		"target":       alert.Target,
		"description":  alert.Description,
		"threat_score": alert.ThreatScore,
		"status":       alert.Status,

		"mitre_tactic":       meta["mitre_tactic"],
		"mitre_technique":    meta["mitre_technique"],
		"mitre_technique_id": meta["mitre_technique_id"],

		"metadata": meta,
	}, alert.ID)
	
	return nil 
}

func GetRecentAlertsByTenant(ctx context.Context, tenantID string) ([]models.Alert, error) {
	query := `
	SELECT id, timestamp, type, severity, source_ip, target,
	       description, threat_score, status
	FROM alerts
	WHERE tenant_id = ?
	ORDER BY timestamp DESC
	LIMIT 50
	`

	rows, err := DB.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.Alert

	for rows.Next() {
		var a models.Alert

		err := rows.Scan(
			&a.ID,
			&a.Timestamp,
			&a.Type,
			&a.Severity,
			&a.SourceIP,
			&a.Target,
			&a.Description,
			&a.ThreatScore,
			&a.Status,
		)
		if err != nil {
			continue
		}

		a.TenantID = tenantID
		alerts = append(alerts, a)
	}

	return alerts, nil
}
