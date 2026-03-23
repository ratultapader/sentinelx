package main

import (
	"fmt"
	"time"

	"sentinelx/response"
)

func main() {
	response.InitResponseEngine(100)
	response.StartActionProcessor()

	response.InitKubernetesExecutor(100)
	k8sController := response.NewKubernetesController(true)
	response.StartKubernetesExecutor(k8sController)

	restartAction := response.Action{
		ID:          "resp_test_restart_1",
		AlertID:     "alert_test_restart_1",
		Timestamp:   time.Now().UTC(),
		ActionType:  response.ActionContainerRestart,
		Severity:    "high",
		ThreatScore: 0.87,
		Reason:      "compromised pod/container suspected",
		Status:      response.StatusPending,
		Metadata: map[string]interface{}{
			"namespace": "payments",
			"pod_name":  "payments-api-7b5c6d7f8d-abc12",
			"node_name": "worker-2",
		},
	}

	isolationAction := response.Action{
		ID:          "resp_test_isolation_1",
		AlertID:     "alert_test_isolation_1",
		Timestamp:   time.Now().UTC(),
		ActionType:  response.ActionK8sIsolation,
		Severity:    "critical",
		ThreatScore: 0.98,
		Reason:      "container/runtime compromise detected",
		Status:      response.StatusPending,
		Metadata: map[string]interface{}{
			"namespace":      "payments",
			"pod_name":       "payments-api-7b5c6d7f8d-abc12",
			"node_name":      "worker-2",
			"isolation_mode": "network_policy",
		},
	}

	fmt.Println("Sending restart action")
	response.ActionQueue <- restartAction

	fmt.Println("Sending isolation action")
	response.ActionQueue <- isolationAction

	time.Sleep(2 * time.Second)
}