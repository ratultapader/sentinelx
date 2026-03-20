package search

import (
	"time"

	"sentinelx/models"
	"sentinelx/storage"
)

func SearchEvents(ip string, eventType string) ([]models.SecurityEvent, error) {
	query := `
	SELECT timestamp, type, source_ip, message
	FROM events
	WHERE 1=1
	`
	args := []interface{}{}

	if ip != "" {
		query += " AND source_ip = ?"
		args = append(args, ip)
	}

	if eventType != "" {
		query += " AND type = ?"
		args = append(args, eventType)
	}

	query += " ORDER BY timestamp DESC"

	rows, err := storage.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]models.SecurityEvent, 0)

	for rows.Next() {
		var ts time.Time
		var typ string
		var sourceIP string
		var message string

		err := rows.Scan(&ts, &typ, &sourceIP, &message)
		if err != nil {
			continue
		}

		event := models.SecurityEvent{
			Timestamp: ts.UnixNano(),
			EventType: typ,
			SourceIP:  sourceIP,
			Metadata: map[string]string{
				"message": message,
			},
		}

		results = append(results, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}