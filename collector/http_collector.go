package collector

import (
	"fmt"      // Used to print the generated event to console
	"net/http" // Provides HTTP server and request handling
	"time"     // Used to measure request processing time
	"sentinelx/pipeline"
	"sentinelx/models"
	"strings"
	"net"

	
)

// HTTPCollector is a middleware.
// It intercepts every HTTP request and creates a security event from it.
func HTTPCollector(next http.Handler) http.Handler {

	// Return a new HTTP handler function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Record the start time of the request
		// Used later to measure latency
		start := time.Now()

		// Call the next handler in the HTTP chain
		// This allows the real application logic to execute
		next.ServeHTTP(w, r)

		// Calculate how long the request took
		duration := time.Since(start)

		// Create a new security event of type "http_request"
		event := models.NewSecurityEvent("http_request")


		// Capture the source IP address of the client
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
event.SourceIP = host
if idx := strings.LastIndex(event.SourceIP, ":"); idx != -1 {
	event.SourceIP = event.SourceIP[:idx]
}

		// Capture the HTTP protocol version (HTTP/1.1, HTTP/2 etc.)
		event.Protocol = r.Proto

		// Store additional request details inside metadata

		// HTTP method used (GET, POST, PUT etc.)
		event.Metadata["method"] = r.Method

		// Requested URL path
		event.Metadata["path"] = r.URL.RequestURI()


		// User agent (browser / client info)
		event.Metadata["user_agent"] = r.UserAgent()

		// How long the request took to complete
		event.Metadata["latency"] = duration.String()

		// If request has body data, store its size
		if r.ContentLength > 0 {
			event.PayloadSize = int(r.ContentLength)
		}

		jsonData, err := event.ToJSON()

if err == nil {
    fmt.Println(string(jsonData))

    // send event to pipeline
	fmt.Println("DEBUG: publishing event to pipeline")
    pipeline.PublishEvent(event)
}

	})
}