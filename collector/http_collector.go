package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"sentinelx/models"
	"sentinelx/multi_tenant"
	"sentinelx/pipeline"
)

// ===============================
// HTTP COLLECTOR (FIXED)
// ===============================
func HTTPCollector(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// ===============================
		// TENANT CHECK
		// ===============================
		tenantID := multi_tenant.TenantIDFromRequest(r)
		if tenantID == "" {
			fmt.Println("no tenant, skipping http event")
			next.ServeHTTP(w, r)
			return
		}

		// ===============================
		// SAFE BODY READ (FIX)
		// ===============================
		var payload struct {
			SourceIP string `json:"source_ip"`
			Payload  string `json:"payload"`
		}

		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)

			// restore body for next handler
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// decode safely
			_ = json.Unmarshal(bodyBytes, &payload)
		}

		// ===============================
		// CONTINUE REQUEST
		// ===============================
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		// ===============================
		// CREATE EVENT
		// ===============================
		event := models.NewSecurityEvent("http_request")
		event.TenantID = tenantID

		// ===============================
		// FIXED SOURCE IP LOGIC
		// ===============================
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		// OVERRIDE WITH PAYLOAD
		if payload.SourceIP != "" {
			host = payload.SourceIP
		}

		event.SourceIP = host

		// ===============================
		// METADATA
		// ===============================
		event.Protocol = r.Proto
		event.Metadata["method"] = r.Method
		event.Metadata["path"] = r.URL.RequestURI()
		event.Metadata["user_agent"] = r.UserAgent()
		event.Metadata["latency"] = duration.String()

		// ===============================
		// ADD PAYLOAD TO METADATA (CRITICAL FIX)
		// ===============================
		if payload.Payload != "" {
			event.Metadata["payload"] = payload.Payload
		}

		if r.ContentLength > 0 {
			event.PayloadSize = int(r.ContentLength)
		}

		// ===============================
		// DEBUG + PIPELINE
		// ===============================
		jsonData, err := event.ToJSON()
		if err == nil {
			fmt.Println(string(jsonData))
			fmt.Println("DEBUG: publishing event to pipeline")

			ctx := multi_tenant.WithTenantID(r.Context(), tenantID)
			pipeline.PublishEvent(ctx, event)
		}
	})
}
