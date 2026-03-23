package response

import (
	"fmt"
	"sync"
)

type KubernetesController struct {
	mu                 sync.Mutex
	simulateMode       bool
	restartedPods      map[string]bool
	isolatedPods       map[string]bool
	cordonedNodes      map[string]bool
	quarantinePolicies map[string]bool
}

func NewKubernetesController(simulate bool) *KubernetesController {
	return &KubernetesController{
		simulateMode:       simulate,
		restartedPods:      make(map[string]bool),
		isolatedPods:       make(map[string]bool),
		cordonedNodes:      make(map[string]bool),
		quarantinePolicies: make(map[string]bool),
	}
}

func (k *KubernetesController) ValidateMetadata(namespace, podName string) error {
	if namespace == "" {
		return fmt.Errorf("missing namespace")
	}
	if podName == "" {
		return fmt.Errorf("missing pod_name")
	}
	return nil
}

func (k *KubernetesController) makePodKey(namespace, podName string) string {
	return namespace + "/" + podName
}

func (k *KubernetesController) RestartPod(namespace, podName string) (string, error) {
	if err := k.ValidateMetadata(namespace, podName); err != nil {
		return "", err
	}

	key := k.makePodKey(namespace, podName)

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.restartedPods[key] {
		return "pod restart already performed", nil
	}

	k.restartedPods[key] = true

	if k.simulateMode {
		return "simulated pod restart executed", nil
	}

	return "pod restart executed", nil
}

func (k *KubernetesController) LabelPodForIsolation(namespace, podName string) (string, error) {
	if err := k.ValidateMetadata(namespace, podName); err != nil {
		return "", err
	}

	key := k.makePodKey(namespace, podName)

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.isolatedPods[key] {
		return "pod already labeled for isolation", nil
	}

	k.isolatedPods[key] = true

	if k.simulateMode {
		return "simulated pod isolation label applied", nil
	}

	return "pod isolation label applied", nil
}

func (k *KubernetesController) CreateQuarantineNetworkPolicy(namespace, podName string) (string, error) {
	if err := k.ValidateMetadata(namespace, podName); err != nil {
		return "", err
	}

	key := k.makePodKey(namespace, podName)

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.quarantinePolicies[key] {
		return "quarantine network policy already exists", nil
	}

	k.quarantinePolicies[key] = true

	if k.simulateMode {
		return "simulated quarantine network policy created", nil
	}

	return "quarantine network policy created", nil
}

func (k *KubernetesController) CordonNode(nodeName string) (string, error) {
	if nodeName == "" {
		return "", fmt.Errorf("missing node_name")
	}

	k.mu.Lock()
	defer k.mu.Unlock()

	if k.cordonedNodes[nodeName] {
		return "node already cordoned", nil
	}

	k.cordonedNodes[nodeName] = true

	if k.simulateMode {
		return "simulated node cordon executed", nil
	}

	return "node cordon executed", nil
}