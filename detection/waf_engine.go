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

	fmt.Println("DEBUG: WAF analyzing request")

	// Only inspect HTTP requests
	if event.EventType != "http_request" {
		return nil
	}

	// Ensure metadata exists
	if event.Metadata == nil {
		return nil
	}

	// ===============================
	// SAFE TYPE CAST (CRITICAL FIX)
	// ===============================
	pathStr := ""
	if v, ok := event.Metadata["path"].(string); ok {
		pathStr = v
	}

	payloadStr := ""
	if v, ok := event.Metadata["payload"].(string); ok {
		payloadStr = v
	}

	// allow detection from either path or payload
	if pathStr == "" && payloadStr == "" {
		return nil
	}

	// ===============================
	// DECODE PATH
	// ===============================
	decodedPath, err := url.QueryUnescape(pathStr)
	if err == nil {
		pathStr = decodedPath
	}

	// COMBINE BOTH (CRITICAL FIX)
	combined := strings.ToLower(pathStr + " " + payloadStr)

	data := combined

	fmt.Println("DEBUG: decoded path =", data)

	// ===============================
	// METADATA
	// ===============================
	metadata := map[string]interface{}{
		"path": pathStr,
	}

	if payloadStr != "" {
		metadata["payload"] = payloadStr
	}

	if method, ok := event.Metadata["method"]; ok {
		metadata["method"] = method
	}

	// ===============================
	// SQL INJECTION
	// ===============================
	if detectPattern(data, sqlInjectionPatterns) {
		return &models.Alert{
			ID:          generateAlertID(),
			TenantID:    event.TenantID,
			Timestamp:   time.Now().UTC(),
			Type:        "sql_injection",
			Severity:    models.SeverityCritical,
			SourceIP:    event.SourceIP,
			Target:      pathStr,
			Description: "SQL injection payload detected",
			ThreatScore: 0.95,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
	}

	// ===============================
	// XSS
	// ===============================
	if detectPattern(data, xssPatterns) {
		return &models.Alert{
			ID:          generateAlertID(),
			TenantID:    event.TenantID,
			Timestamp:   time.Now().UTC(),
			Type:        "xss_attack",
			Severity:    models.SeverityHigh,
			SourceIP:    event.SourceIP,
			Target:      pathStr,
			Description: "Cross-site scripting payload detected",
			ThreatScore: 0.80,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
	}

	// ===============================
	// DIRECTORY TRAVERSAL
	// ===============================
	if detectPattern(data, traversalPatterns) {
		return &models.Alert{
			ID:          generateAlertID(),
			TenantID:    event.TenantID,
			Timestamp:   time.Now().UTC(),
			Type:        "dir_traversal",
			Severity:    models.SeverityHigh,
			SourceIP:    event.SourceIP,
			Target:      pathStr,
			Description: "Directory traversal payload detected",
			ThreatScore: 0.85,
			Status:      models.AlertStatusNew,
			Metadata:    metadata,
		}
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
