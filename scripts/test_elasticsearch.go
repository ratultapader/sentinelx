package main

import (
	"context"
	"fmt"
	"time"

	"sentinelx/storage"
)

func main() {
	fmt.Println("DEBUG 1: starting test")

	cfg := storage.ElasticsearchConfig{
		Enabled:   true,
		Addresses: []string{"http://localhost:9200"},
		Username:  "",
		Password:  "",
	}

	fmt.Println("DEBUG 2: before InitElasticsearch")
	err := storage.InitElasticsearch(cfg)
	fmt.Println("DEBUG 3: after InitElasticsearch, err =", err)
	if err != nil {
		panic(err)
	}

	fmt.Println("DEBUG 4: before IndexAlertDoc")
	storage.IndexAlertDoc(map[string]interface{}{
		"id":           "alert_test_es_1",
		"timestamp":    time.Now().UTC(),
		"type":         "xss_attack",
		"severity":     "high",
		"source_ip":    "192.168.1.50",
		"target":       "/search?q=<script>alert(1)</script>",
		"description":  "Cross-site scripting payload detected",
		"threat_score": 0.82,
		"status":       "new",
		"metadata": map[string]interface{}{
			"method": "GET",
			"path":   "/search?q=<script>alert(1)</script>",
		},
	}, "alert_test_es_1")
	fmt.Println("DEBUG 5: after IndexAlertDoc")

	fmt.Println("DEBUG 6: before IndexResponseActionDoc")
	storage.IndexResponseActionDoc(map[string]interface{}{
		"id":           "resp_test_es_1",
		"alert_id":     "alert_test_es_1",
		"timestamp":    time.Now().UTC(),
		"action_type":  "rate_limit",
		"source_ip":    "192.168.1.50",
		"target":       "/search?q=<script>alert(1)</script>",
		"severity":     "high",
		"threat_score": 0.82,
		"reason":       "high threat score or severity",
		"status":       "pending",
		"metadata": map[string]interface{}{
			"alert_type": "xss_attack",
		},
	}, "resp_test_es_1")
	fmt.Println("DEBUG 7: after IndexResponseActionDoc")

	time.Sleep(2 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("DEBUG 8: before SearchBySourceIP alerts")
	alertResults, err := storage.ESStore.SearchBySourceIP(ctx, storage.IndexAlerts, "192.168.1.50")
	fmt.Println("DEBUG 9: after SearchBySourceIP alerts, err =", err)
	if err != nil {
		panic(err)
	}

	fmt.Println("DEBUG 10: before SearchBySourceIP response_actions")
	actionResults, err := storage.ESStore.SearchBySourceIP(ctx, storage.IndexResponseActions, "192.168.1.50")
	fmt.Println("DEBUG 11: after SearchBySourceIP response_actions, err =", err)
	if err != nil {
		panic(err)
	}

	fmt.Println("Alert documents found:", len(alertResults))
	fmt.Println("Response action documents found:", len(actionResults))
}