package main

import (
	"fmt"

	"sentinelx/investigation"
)

func main() {
	alertDoc := map[string]interface{}{
		"id":           "alert_1",
		"timestamp":    "2026-03-20T10:02:00Z",
		"source_ip":    "192.168.1.5",
		"event_type":   "sql_injection",
		"severity":     "critical",
		"threat_score": 0.91,
		"description":  "SQL injection payload detected on /login",
	}

	eventDoc := map[string]interface{}{
		"id":           "event_1",
		"timestamp":    "2026-03-20T10:00:00Z",
		"source_ip":    "192.168.1.5",
		"event_type":   "port_scan",
		"severity":     "medium",
		"threat_score": 0.55,
		"message":      "Multiple port probes detected",
	}

	actionDoc := map[string]interface{}{
		"id":           "action_1",
		"timestamp":    "2026-03-20T10:07:00Z",
		"source_ip":    "192.168.1.5",
		"action_type":  "ip_block",
		"severity":     "critical",
		"threat_score": 0.91,
		"reason":       "threat score exceeded ip-block threshold",
	}

	resultDoc := map[string]interface{}{
		"id":           "result_1",
		"timestamp":    "2026-03-20T10:07:01Z",
		"source_ip":    "192.168.1.5",
		"event":        "ip_blocked",
		"severity":     "critical",
		"threat_score": 0.91,
		"message":      "IP successfully blocked",
	}

	events := []investigation.TimelineEvent{
		investigation.NormalizeSecurityEvent(eventDoc),
		investigation.NormalizeAlert(alertDoc),
		investigation.NormalizeResponseAction(actionDoc),
		investigation.NormalizeResponseResult(resultDoc),
	}

	builder := investigation.NewBuilder()
	timeline := builder.Build("192.168.1.5", events)

	fmt.Println(investigation.GenerateTextReport(timeline))
}