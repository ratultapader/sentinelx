package incident

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"sentinelx/models"
)

var (
	incidents = make(map[string]models.Incident)
	mutex     sync.Mutex
)

// CreateIncidentFromAlert creates an incident directly from a full alert model.
func CreateIncidentFromAlert(alert models.Alert) models.Incident {
	mutex.Lock()
	defer mutex.Unlock()

	id := uuid.New().String()

	description := alert.Description
	if description == "" {
		description = "Incident created from security alert"
	}

	incident := models.Incident{
		ID:          id,
		Title:       alert.Type,
		Description: description,
		Severity:    alert.Severity,
		Status:      "open",
		CreatedAt:   time.Now(),
		Alerts:      []string{alert.ID},
	}

	incidents[id] = incident
	return incident
}

func GetAllIncidents() []models.Incident {
	mutex.Lock()
	defer mutex.Unlock()

	list := make([]models.Incident, 0, len(incidents))

	for _, inc := range incidents {
		list = append(list, inc)
	}

	return list
}