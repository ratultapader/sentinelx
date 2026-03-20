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

func CreateIncident(alertID string, title string, severity string) models.Incident {
	mutex.Lock()
	defer mutex.Unlock()

	id := uuid.New().String()

	incident := models.Incident{
		ID:          id,
		Title:       title,
		Description: "Incident created from security alert",
		Severity:    severity,
		Status:      "open",
		CreatedAt:   time.Now(),
		Alerts:      []string{alertID},
	}

	incidents[id] = incident
	return incident
}

func GetAllIncidents() []models.Incident {
	mutex.Lock()
	defer mutex.Unlock()

	list := make([]models.Incident, 0)

	for _, inc := range incidents {
		list = append(list, inc)
	}

	return list
}