package metrics

import (
	"fmt"
	"time"
)

func StartMetricsReporter() {

	ticker := time.NewTicker(10 * time.Second)

	for range ticker.C {
		PrintMetrics()
	}
}

func PrintMetrics() {

	mu.RLock()
	defer mu.RUnlock()

	fmt.Println("========== SentinelX Metrics ==========")

	fmt.Println("Total Events:", TotalEvents)
	fmt.Println("Total Alerts:", TotalAlerts)

	fmt.Println("Top Attack Types:")
	for k, v := range AttackTypes {
		fmt.Printf("   %s : %d\n", k, v)
	}

	fmt.Println("Top Attacker IPs:")
	for k, v := range AttackerIPs {
		fmt.Printf("   %s : %d\n", k, v)
	}

	fmt.Println("=======================================")
}