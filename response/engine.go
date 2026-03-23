package response

import (
	"fmt"
	"strings"
	"time"

	"sentinelx/models"
)

func Decide(alert models.Alert) Action {
	alert = normalizeThreatScore(alert)
	actionType, reason := determineAction(alert)

	return Action{
		ID:          generateResponseID(alert.ID),
		AlertID:     alert.ID,
		Timestamp:   time.Now().UTC(),
		ActionType:  actionType,
		SourceIP:    alert.SourceIP,
		Target:      alert.Target,
		Severity:    alert.Severity,
		ThreatScore: alert.ThreatScore,
		Reason:      reason,
		Status:      StatusPending,
		Metadata: map[string]interface{}{
			"alert_type": alert.Type,
		},
	}
}

func determineAction(alert models.Alert) (string, string) {
	alertType := strings.ToLower(alert.Type)

	switch {
	case alertType == "container_escape":
		return ActionK8sIsolation, "container escape requires kubernetes isolation"
	case alertType == "runtime_compromise" || alertType == "malicious_container_activity":
		return ActionK8sIsolation, "container/runtime compromise detected"
	case alertType == "compromised_pod" || alertType == "container_backdoor":
		return ActionContainerRestart, "compromised pod/container suspected"
	case alert.ThreatScore > 0.9:
		return ActionIPBlock, "threat score exceeded ip-block threshold"
	case alert.ThreatScore > 0.7:
		return ActionRateLimit, "threat score exceeded rate-limit threshold"
	case alert.ThreatScore > 0.5:
		return ActionAlertOnly, "threat score exceeded alert-only threshold"
	default:
		return ActionAlertOnly, "monitoring only; below mitigation threshold"
	}
}

func normalizeThreatScore(alert models.Alert) models.Alert {
	if alert.ThreatScore > 0 {
		return alert
	}

	switch strings.ToLower(alert.Severity) {
	case "critical":
		alert.ThreatScore = 0.95
	case "high":
		alert.ThreatScore = 0.80
	case "medium":
		alert.ThreatScore = 0.60
	case "low":
		alert.ThreatScore = 0.30
	default:
		alert.ThreatScore = 0.10
	}

	return alert
}

func generateResponseID(alertID string) string {
	return fmt.Sprintf("resp_%s_%d", alertID, time.Now().UnixNano())
}
