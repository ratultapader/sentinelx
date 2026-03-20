package storage

import (
	"fmt"
	"time"

	"sentinelx/models"
)

func SaveEvent(event models.SecurityEvent) {
	query := `
	INSERT INTO events (timestamp, type, source_ip, message)
	VALUES (?, ?, ?, ?)
	`

	message := fmt.Sprintf("event_type=%s", event.EventType)

	if path, ok := event.Metadata["path"]; ok {
		message = path
	} else if method, ok := event.Metadata["method"]; ok {
		message = fmt.Sprintf("%s %s", method, event.EventType)
	}

	DB.Exec(
		query,
		time.Now(),
		event.EventType,
		event.SourceIP,
		message,
	)
}