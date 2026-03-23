package response

import (
	"fmt"
	"time"
)

var FirewallResults chan FirewallResult

func InitFirewallExecutor(size int) {
	FirewallResults = make(chan FirewallResult, size)
	fmt.Println("Firewall Executor initialized")
}

func StartFirewallExecutor(blocker *FirewallBlocker) {
	fmt.Println("Firewall Executor started")

	go func() {
		for action := range FirewallActionQueue {
			if action.ActionType != ActionIPBlock {
				logFirewallResult(FirewallResult{
					ID:         generateFirewallResultID(action.ID),
					ActionID:   action.ID,
					AlertID:    action.AlertID,
					Timestamp:  time.Now().UTC(),
					Event:      "firewall_skipped",
					SourceIP:   action.SourceIP,
					ActionType: action.ActionType,
					Status:     "skipped",
					Message:    "action is not ip_block",
					Metadata: map[string]interface{}{
						"reason": action.Reason,
					},
				})
				continue
			}

			result := FirewallResult{
				ID:         generateFirewallResultID(action.ID),
				ActionID:   action.ID,
				AlertID:    action.AlertID,
				Timestamp:  time.Now().UTC(),
				Event:      "ip_blocked",
				SourceIP:   action.SourceIP,
				ActionType: action.ActionType,
				Status:     StatusPending,
				Message:    "processing firewall block request",
				Metadata: map[string]interface{}{
					"severity":     action.Severity,
					"threat_score": action.ThreatScore,
					"reason":       action.Reason,
				},
			}

			msg, err := blocker.BlockIP(action.SourceIP)
			if err != nil {
				result.Status = StatusFailed
				result.Message = err.Error()
				logFirewallResult(result)
				continue
			}

			switch msg {
			case "source ip is protected and cannot be blocked automatically":
				result.Status = "skipped"
				result.Message = msg
			case "ip already blocked":
				result.Status = "skipped"
				result.Message = msg
			case "simulated firewall block applied", "firewall block applied":
				result.Status = StatusExecuted
				result.Message = msg
			default:
				result.Status = StatusExecuted
				result.Message = msg
			}

			logFirewallResult(result)
		}
	}()
}

func logFirewallResult(result FirewallResult) {
	fmt.Println("========= FIREWALL RESULT =========")
	fmt.Println("ID:", result.ID)
	fmt.Println("Action ID:", result.ActionID)
	fmt.Println("Alert ID:", result.AlertID)
	fmt.Println("Event:", result.Event)
	fmt.Println("Source IP:", result.SourceIP)
	fmt.Println("Action Type:", result.ActionType)
	fmt.Println("Status:", result.Status)
	fmt.Println("Message:", result.Message)
	fmt.Println("Timestamp:", result.Timestamp)
	fmt.Println("===================================")

	if FirewallResults != nil {
		select {
		case FirewallResults <- result:
		default:
			fmt.Println("Firewall result queue full — dropping result")
		}
	}
}

func generateFirewallResultID(actionID string) string {
	return fmt.Sprintf("fw_%s_%d", actionID, time.Now().UnixNano())
}