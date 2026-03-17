package storage

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
	ID        int    `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Severity  string `json:"severity"`
	SourceIP  string `json:"source_ip"`
	Message   string `json:"message"`
}

func GetAlerts(severity string) ([]AlertRecord, error) {

	// ⚠ IMPORTANT: your DB column is "description"
	query := "SELECT id, timestamp, type, severity, source_ip, description FROM alerts WHERE 1=1"

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

		err := rows.Scan(
			&a.ID,
			&a.Timestamp,
			&a.Type,
			&a.Severity,
			&a.SourceIP,
			&a.Message, // maps description → message
		)

		if err != nil {
			return nil, err
		}

		alerts = append(alerts, a)
	}

	return alerts, nil
}