package detection

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"sentinelx/models"
)

type WAFEngine struct{}

var WAF = &WAFEngine{}

// SQL Injection signatures
var sqlInjectionPatterns = []string{
	"or 1=1",
	"' or 1=1",
	"union select",
	"drop table",
	"select * from",
	"'--",
	"1=1",
}

// XSS signatures
var xssPatterns = []string{
	"<script>",
	"</script>",
	"javascript:",
}

// Directory traversal signatures
var traversalPatterns = []string{
	"../",
	"..\\",
}

// ProcessEvent analyzes HTTP request events and returns generated alerts.
func (w *WAFEngine) ProcessEvent(event models.SecurityEvent) *models.Alert {
	// Only inspect HTTP requests
	fmt.Println("DEBUG: WAF analyzing request")
	if event.EventType != "http_request" {
		return nil
	}

	// Ensure metadata exists
	if event.Metadata == nil {
		return nil
	}

	path, exists := event.Metadata["path"]
	if !exists {
		return nil
	}

	decodedPath, err := url.QueryUnescape(path)
	if err == nil {
		path = decodedPath
	}

	data := strings.ToLower(path)

	// DEBUG point
	fmt.Println("DEBUG: decoded path =", data)

	metadata := map[string]interface{}{
		"path": path,
	}

	if method, ok := event.Metadata["method"]; ok {
		metadata["method"] = method
	}

	// SQL Injection detection
	if detectPattern(data, sqlInjectionPatterns) {
		alert := models.Alert{
			ID:          generateAlertID(),
			Timestamp:   time.Now().UTC(),
			Type:        "sql_injection",
			Severity:    models.SeverityCritical,
			SourceIP:    event.SourceIP,
			Target:      path,
			Description: "SQL injection payload detected",
			ThreatScore: 0.95,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
		return &alert
	}

	// XSS detection
	if detectPattern(data, xssPatterns) {
		alert := models.Alert{
			ID:          generateAlertID(),
			Timestamp:   time.Now().UTC(),
			Type:        "xss_attack",
			Severity:    models.SeverityHigh,
			SourceIP:    event.SourceIP,
			Target:      path,
			Description: "Cross-site scripting payload detected",
			ThreatScore: 0.80,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
		return &alert
	}

	// Directory traversal detection
	if detectPattern(data, traversalPatterns) {
		alert := models.Alert{
			ID:          generateAlertID(),
			Timestamp:   time.Now().UTC(),
			Type:        "dir_traversal",
			Severity:    models.SeverityHigh,
			SourceIP:    event.SourceIP,
			Target:      path,
			Description: "Directory traversal payload detected",
			ThreatScore: 0.85,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
		return &alert
	}

	return nil
}

// detectPattern checks if any signature exists.
func detectPattern(data string, patterns []string) bool {
	for _, pattern := range patterns {
		fmt.Println("DEBUG: checking pattern:", pattern)

		if strings.Contains(data, pattern) {
			fmt.Println("DEBUG: pattern matched:", pattern)
			return true
		}
	}

	return false
}