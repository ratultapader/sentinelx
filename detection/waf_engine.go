package detection

import (
	 "fmt"
"net/url"
	"strings"


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

// ProcessEvent analyzes HTTP request events
func (w *WAFEngine) ProcessEvent(event models.SecurityEvent) {

	// Only inspect HTTP requests
	fmt.Println("DEBUG: WAF analyzing request")
	if event.EventType != "http_request" {
		return
	}

	// Ensure metadata exists
	if event.Metadata == nil {
		return
	}

	path, exists := event.Metadata["path"]
if !exists {
	return
}

decodedPath, err := url.QueryUnescape(path)
if err == nil {
	path = decodedPath
}

data := strings.ToLower(path)

// DEBUG point
fmt.Println("DEBUG: decoded path =", data)

	// SQL Injection detection
	if detectPattern(data, sqlInjectionPatterns) {

		GenerateAlert(
    "SQL_INJECTION",
    event.SourceIP,
    "SQL injection attempt detected",
)

		return
	}

	// XSS detection
	if detectPattern(data, xssPatterns) {

		GenerateAlert(
    "XSS_ATTACK",
    event.SourceIP,
    "Cross-site scripting attempt detected",
)

		return
	}

	// Directory traversal detection
	if detectPattern(data, traversalPatterns) {

		GenerateAlert(
    "DIR_TRAVERSAL",
    event.SourceIP,
    "Directory traversal attempt detected",
)

		return
	}
}

// detectPattern checks if any signature exists
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