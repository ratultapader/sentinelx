package detection

import (
	"fmt"
	"time"

	"sentinelx/models"
	"sentinelx/storage"
)

// Global alert queue
var AlertQueue chan models.Alert

// InitAlertEngine initializes the alert system
func InitAlertEngine(size int) {

	AlertQueue = make(chan models.Alert, size)

}

// GenerateAlert creates and sends alert into queue
func GenerateAlert(severity, alertType, ip, description string) {

	alert := models.Alert{
		Timestamp:   time.Now().Unix(),
		Severity:    severity,
		Type:        alertType,
		SourceIP:    ip,
		Description: description,
	}

	select {

	case AlertQueue <- alert:

	default:
		// queue full → drop alert
	}

}

// StartAlertProcessor consumes alerts
func StartAlertProcessor() {

	for alert := range AlertQueue {

		fmt.Println("SECURITY ALERT")
		fmt.Println("Type:", alert.Type)
		fmt.Println("Source IP:", alert.SourceIP)
		fmt.Println("Severity:", alert.Severity)
		fmt.Println("Description:", alert.Description)
		fmt.Println("-------------")

		storage.LogAlert(alert)

	}

}
