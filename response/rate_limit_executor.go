package response

import (
	"fmt"
	"time"
)

var RateLimitResults chan RateLimitResult

func InitRateLimitExecutor(size int) {
	RateLimitResults = make(chan RateLimitResult, size)
	fmt.Println("Rate Limit Executor initialized")
}

func StartRateLimitExecutor(limiter *RateLimiter, limitPerSec int, burst int) {
	fmt.Println("Rate Limit Executor started")

	go func() {
		for action := range RateLimitActionQueue {
			result := RateLimitResult{
				ID:          generateRateLimitResultID(action.ID),
				ActionID:    action.ID,
				AlertID:     action.AlertID,
				Timestamp:   time.Now().UTC(),
				Event:       "rate_limit_applied",
				SourceIP:    action.SourceIP,
				ActionType:  action.ActionType,
				Status:      StatusPending,
				Message:     "processing rate limit request",
				LimitPerSec: limitPerSec,
				Burst:       burst,
				Metadata: map[string]interface{}{
					"severity":     action.Severity,
					"threat_score": action.ThreatScore,
					"reason":       action.Reason,
				},
			}

			msg, err := limiter.ApplyLimit(action.SourceIP, limitPerSec, burst)
			if err != nil {
				result.Status = StatusFailed
				result.Message = err.Error()
				logRateLimitResult(result)
				continue
			}

			switch msg {
			case "source ip is protected and cannot be rate limited automatically":
				result.Status = "skipped"
				result.Message = msg
			case "rate limit already exists for source ip":
				result.Status = "skipped"
				result.Message = msg
			case "simulated rate limit applied", "rate limit applied":
				result.Status = StatusExecuted
				result.Message = msg
			default:
				result.Status = StatusExecuted
				result.Message = msg
			}

			logRateLimitResult(result)
		}
	}()
}

func logRateLimitResult(result RateLimitResult) {
	fmt.Println("======= RATE LIMIT RESULT =========")
	fmt.Println("ID:", result.ID)
	fmt.Println("Action ID:", result.ActionID)
	fmt.Println("Alert ID:", result.AlertID)
	fmt.Println("Event:", result.Event)
	fmt.Println("Source IP:", result.SourceIP)
	fmt.Println("Action Type:", result.ActionType)
	fmt.Println("Status:", result.Status)
	fmt.Println("Message:", result.Message)
	fmt.Println("Limit/sec:", result.LimitPerSec)
	fmt.Println("Burst:", result.Burst)
	fmt.Println("Timestamp:", result.Timestamp)
	fmt.Println("==================================")

	if RateLimitResults != nil {
		select {
		case RateLimitResults <- result:
		default:
			fmt.Println("Rate-limit result queue full — dropping result")
		}
	}
}

func generateRateLimitResultID(actionID string) string {
	return fmt.Sprintf("rl_%s_%d", actionID, time.Now().UnixNano())
}