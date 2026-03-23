package response

import (
	"fmt"
	"time"
)

var KubernetesResults chan KubernetesResult

func InitKubernetesExecutor(size int) {
	KubernetesResults = make(chan KubernetesResult, size)
	fmt.Println("Kubernetes Executor initialized")
}

func StartKubernetesExecutor(controller *KubernetesController) {
	fmt.Println("Kubernetes Executor started")

	go func() {
		for action := range KubernetesActionQueue {
			namespace := getMetadataString(action.Metadata, "namespace")
			podName := getMetadataString(action.Metadata, "pod_name")
			nodeName := getMetadataString(action.Metadata, "node_name")
			isolationMode := getMetadataString(action.Metadata, "isolation_mode")

			result := KubernetesResult{
				ID:         generateKubernetesResultID(action.ID),
				ActionID:   action.ID,
				AlertID:    action.AlertID,
				Timestamp:  time.Now().UTC(),
				Event:      "kubernetes_action_executed",
				ActionType: action.ActionType,
				Namespace:  namespace,
				PodName:    podName,
				NodeName:   nodeName,
				Status:     StatusPending,
				Message:    "processing kubernetes action",
				Metadata: map[string]interface{}{
					"severity":       action.Severity,
					"threat_score":   action.ThreatScore,
					"reason":         action.Reason,
					"isolation_mode": isolationMode,
				},
			}

			switch action.ActionType {
			case ActionContainerRestart:
				msg, err := controller.RestartPod(namespace, podName)
				if err != nil {
					result.Status = StatusFailed
					result.Message = err.Error()
					logKubernetesResult(result)
					continue
				}

				switch msg {
				case "pod restart already performed":
					result.Status = "skipped"
					result.Event = "pod_restart_skipped"
				default:
					result.Status = StatusExecuted
					result.Event = "pod_restarted"
				}
				result.Message = msg

			case ActionK8sIsolation:
				labelMsg, err := controller.LabelPodForIsolation(namespace, podName)
				if err != nil {
					result.Status = StatusFailed
					result.Message = err.Error()
					logKubernetesResult(result)
					continue
				}

				policyMsg, err := controller.CreateQuarantineNetworkPolicy(namespace, podName)
				if err != nil {
					result.Status = StatusFailed
					result.Message = err.Error()
					logKubernetesResult(result)
					continue
				}

				finalMessage := labelMsg + "; " + policyMsg

				if nodeName != "" {
					cordonMsg, err := controller.CordonNode(nodeName)
					if err != nil {
						result.Status = StatusFailed
						result.Message = err.Error()
						logKubernetesResult(result)
						continue
					}
					finalMessage += "; " + cordonMsg
				}

				if labelMsg == "pod already labeled for isolation" &&
					policyMsg == "quarantine network policy already exists" {
					result.Status = "skipped"
					result.Event = "kubernetes_isolation_skipped"
				} else {
					result.Status = StatusExecuted
					result.Event = "kubernetes_isolated"
				}

				result.Message = finalMessage

			default:
				result.Status = "skipped"
				result.Event = "kubernetes_action_skipped"
				result.Message = "unsupported kubernetes action type"
			}

			logKubernetesResult(result)
		}
	}()
}

func logKubernetesResult(result KubernetesResult) {
	fmt.Println("======= KUBERNETES RESULT =========")
	fmt.Println("ID:", result.ID)
	fmt.Println("Action ID:", result.ActionID)
	fmt.Println("Alert ID:", result.AlertID)
	fmt.Println("Event:", result.Event)
	fmt.Println("Action Type:", result.ActionType)
	fmt.Println("Namespace:", result.Namespace)
	fmt.Println("Pod Name:", result.PodName)
	fmt.Println("Node Name:", result.NodeName)
	fmt.Println("Status:", result.Status)
	fmt.Println("Message:", result.Message)
	fmt.Println("Timestamp:", result.Timestamp)
	fmt.Println("===================================")

	if KubernetesResults != nil {
		select {
		case KubernetesResults <- result:
		default:
			fmt.Println("Kubernetes result queue full — dropping result")
		}
	}
}

func generateKubernetesResultID(actionID string) string {
	return fmt.Sprintf("k8s_%s_%d", actionID, time.Now().UnixNano())
}

func getMetadataString(metadata map[string]interface{}, key string) string {
	if metadata == nil {
		return ""
	}

	v, ok := metadata[key]
	if !ok || v == nil {
		return ""
	}

	s, ok := v.(string)
	if !ok {
		return ""
	}

	return s
}