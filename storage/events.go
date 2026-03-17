package storage

import (
	"time"

	"sentinelx/models"
)

func SaveEvent(event models.SecurityEvent) {

	query := `
	INSERT INTO events (timestamp, type, source_ip, message)
	VALUES (?, ?, ?, ?)
	`

	DB.Exec(
		query,
		time.Now(),
		event.EventType,
		event.SourceIP,
		event.EventType, // or message if you have one
	)
}