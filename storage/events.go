package storage

import (
	"fmt"
	"log"
	"time"
	"context"

	"sentinelx/models"
)

func SaveEvent(event models.SecurityEvent) {
	query := `
	INSERT INTO events (timestamp, type, source_ip, message)
	VALUES (?, ?, ?, ?)
	`

	message := fmt.Sprintf("event_type=%s", event.EventType)

	if path, ok := event.Metadata["path"]; ok {
		message = fmt.Sprintf("%v", path)
	} else if method, ok := event.Metadata["method"]; ok {
		message = fmt.Sprintf("%v %s", method, event.EventType)
	}

	_, err := DB.Exec(
		query,
		time.Now(),
		event.EventType,
		event.SourceIP,
		message,
	)
	if err != nil {
		log.Println("failed to save event:", err)
		return
	}

	doc := map[string]interface{}{
	"id":         event.EventID,
	"tenant_id":  event.TenantID, // ✅ ADD THIS LINE
	"timestamp":  time.Unix(0, event.Timestamp).UTC(),
	"event_type": event.EventType,
	"source_ip":  event.SourceIP,
	"protocol":   event.Protocol,
	"metadata":   event.Metadata,
}

IndexSecurityEventDoc(context.Background(), doc, event.EventID)
}