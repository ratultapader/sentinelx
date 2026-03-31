package response

import (
	"context"
	"fmt"

	"sentinelx/models"
	"sentinelx/storage"
)

var ActionQueue chan Action
var FirewallActionQueue chan Action
var RateLimitActionQueue chan Action
var KubernetesActionQueue chan Action

func InitResponseEngine(size int) {
	ActionQueue = make(chan Action, size)
	FirewallActionQueue = make(chan Action, size)
	RateLimitActionQueue = make(chan Action, size)
	KubernetesActionQueue = make(chan Action, size)
	fmt.Println("Response Engine initialized")
}

func ProcessAlert(alert models.Alert) {
	action := Decide(alert)
	action.TenantID = alert.TenantID

	if storage.GraphIngestor != nil {
		err := storage.GraphIngestor.IngestAttackRecord(context.Background(), storage.AttackGraphRecord{
			TenantID:         alert.TenantID,
			AlertID:          alert.ID,
			Timestamp:        alert.Timestamp,
			SourceIP:         alert.SourceIP,
			Server:           alert.Target,
			APIEndpoint:      alert.Target,
			EventType:        alert.Type,
			Severity:         alert.Severity,
			ThreatScore:      alert.ThreatScore,
			ResponseAction:   action.ActionType,
			MitreTactic:      fmt.Sprintf("%v", alert.Metadata["mitre_tactic"]),
			MitreTechnique:   fmt.Sprintf("%v", alert.Metadata["mitre_technique"]),
			MitreTechniqueID: fmt.Sprintf("%v", alert.Metadata["mitre_technique_id"]),
		})
		if err != nil {
			fmt.Println("WARNING: Neo4j graph ingest failed:", err)
		}
	}

	logResponseAction(action)

	select {
	case ActionQueue <- action:
	default:
		fmt.Println("Response action queue full - dropping action")
	}
}

func StartActionProcessor() {
	fmt.Println("Response Action Processor started")

	go func() {
		for action := range ActionQueue {
			switch action.ActionType {
			case ActionAlertOnly:
				fmt.Println("ALERT ONLY -> no action needed")
			case ActionIPBlock:
				select {
				case FirewallActionQueue <- action:
				default:
					fmt.Println("Firewall action queue full - dropping action")
				}
			case ActionRateLimit:
				select {
				case RateLimitActionQueue <- action:
				default:
					fmt.Println("Rate-limit action queue full - dropping action")
				}
			case ActionContainerRestart, ActionK8sIsolation:
				select {
				case KubernetesActionQueue <- action:
				default:
					fmt.Println("Kubernetes action queue full - dropping action")
				}
			default:
				fmt.Println("Unknown action:", action.ActionType)
			}
		}
	}()
}

func logResponseAction(action Action) {
	fmt.Println("========= RESPONSE ACTION =========")
	fmt.Println("ID:", action.ID)
	fmt.Println("Alert ID:", action.AlertID)
	fmt.Println("Action:", action.ActionType)
	fmt.Println("Source IP:", action.SourceIP)
	fmt.Println("Target:", action.Target)
	fmt.Println("Severity:", action.Severity)
	fmt.Println("Threat Score:", action.ThreatScore)
	fmt.Println("Reason:", action.Reason)
	fmt.Println("Status:", action.Status)
	fmt.Println("Timestamp:", action.Timestamp)
	fmt.Println("===================================")

	ctx := context.Background()

	storage.IndexResponseActionDoc(ctx, map[string]interface{}{
		"id":           action.ID,
		"alert_id":     action.AlertID,
		"timestamp":    action.Timestamp,
		"action_type":  action.ActionType,
		"source_ip":    action.SourceIP,
		"target":       action.Target,
		"severity":     action.Severity,
		"threat_score": action.ThreatScore,
		"reason":       action.Reason,
		"status":       action.Status,
		"metadata":     action.Metadata,
		"tenant_id":    action.TenantID,
	}, action.ID)
}


