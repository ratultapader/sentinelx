package storage

import (
	"time"

	"sentinelx/models"
)

func SaveAlert(alert models.Alert) {

	query := `
	INSERT INTO alerts (timestamp, type, severity, source_ip, description)
	VALUES (?, ?, ?, ?, ?)
	`

	DB.Exec(
		query,
		time.Now(),
		alert.Type,
		alert.Severity,
		alert.SourceIP,
		alert.Description,
	)
}