package storage

import "encoding/json"

type EventRecord struct {
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	SourceIP  string `json:"source_ip"`
	Message   string `json:"message"`
}

func GetEvents(ip string, eventType string) ([]EventRecord, error) {
	query := "SELECT id, timestamp, type, source_ip, message FROM events WHERE 1=1"
	args := []interface{}{}

	if ip != "" {
		query += " AND source_ip=?"
		args = append(args, ip)
	}

	if eventType != "" {
		query += " AND type=?"
		args = append(args, eventType)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []EventRecord

	for rows.Next() {
		var e EventRecord

		err := rows.Scan(
			&e.ID,
			&e.Timestamp,
			&e.Type,
			&e.SourceIP,
			&e.Message,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

// ================= ALERTS =================

type AlertRecord struct {
	ID          string                 `json:"id"`
	Timestamp   string                 `json:"timestamp"`
	Type        string                 `json:"type"`
	Severity    string                 `json:"severity"`
	SourceIP    string                 `json:"source_ip"`
	Target      string                 `json:"target,omitempty"`
	Description string                 `json:"description"`
	ThreatScore float64                `json:"threat_score"`
	Status      string                 `json:"status"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

func GetAlerts(severity string) ([]AlertRecord, error) {
	query := `
	SELECT id, timestamp, type, severity, source_ip, target, description, threat_score, status, metadata
	FROM alerts
	WHERE 1=1
	`
	args := []interface{}{}

	if severity != "" {
		query += " AND severity=?"
		args = append(args, severity)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []AlertRecord

	for rows.Next() {
		var a AlertRecord
		var metadataJSON []byte

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
			&metadataJSON,
		)
		if err != nil {
			return nil, err
		}

		if len(metadataJSON) > 0 {
			err := json.Unmarshal(metadataJSON, &a.Metadata)
			if err != nil {
				a.Metadata = map[string]interface{}{}
			}
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}