package main

import (
	"fmt"
	"time"

	"sentinelx/api"
	"sentinelx/collector"
	"sentinelx/correlation"
	"sentinelx/detection"
	"sentinelx/metrics"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/response"
	"sentinelx/ruleengine"
	"sentinelx/storage"
	"sentinelx/threatfeed"
	"sentinelx/threatintel"
)

const (
	EventQueueSize = 10000
	AlertQueueSize = 1000
	WorkerCount    = 5
)

func main() {
	fmt.Println("Starting SentinelX Security Platform")

	// Load threat intelligence feed
	err := threatintel.LoadThreatFeed("data/malicious_ips.txt")
	if err != nil {
		panic(err)
	}

	// Initialize logger
	fmt.Println("Initializing event logger...")
	err = storage.InitLogger("logs/security_events.json")
	if err != nil {
		panic(err)
	}

	// Initialize database
	err = storage.InitDB()
	if err != nil {
		panic(err)
	}

	// Initialize Elasticsearch forensic store
	err = storage.InitElasticsearch(storage.ElasticsearchConfig{
		Enabled:   true,
		Addresses: []string{"http://localhost:9200"},
		Username:  "",
		Password:  "",
	})
	if err != nil {
		fmt.Println("WARNING: Elasticsearch init failed:", err)
	} else {
		fmt.Println("Elasticsearch forensic store initialized")
	}

	// Initialize Neo4j graph engine
	storage.InitNeo4jGraph("neo4j://localhost:7687", "neo4j", "password", "neo4j")
	defer storage.CloseNeo4jGraph()

	// Load detection rules
	fmt.Println("Loading detection rules...")
	err = ruleengine.LoadRules()
	if err != nil {
		panic(err)
	}

	// Initialize central alert engine
	fmt.Println("Initializing alert engine...")
	detection.InitAlertEngine(AlertQueueSize)
	go detection.StartAlertProcessor()

	// Initialize response engine
	fmt.Println("Initializing response engine...")
	response.InitResponseEngine(1000)
	go response.StartActionProcessor()

	fmt.Println("Initializing firewall executor...")
	response.InitFirewallExecutor(1000)
	blocker := response.NewFirewallBlocker(true) // true = simulate mode
	response.StartFirewallExecutor(blocker)

	fmt.Println("Initializing rate limit executor...")
	response.InitRateLimitExecutor(1000)
	limiter := response.NewRateLimiter(true)
	response.StartRateLimitExecutor(limiter, 20, 40)

	fmt.Println("Initializing kubernetes executor...")
	response.InitKubernetesExecutor(1000)
	k8sController := response.NewKubernetesController(true)
	response.StartKubernetesExecutor(k8sController)

	// Start threat feed updater
	fmt.Println("Starting threat feed updater...")
	threatfeed.StartThreatFeedUpdater()

	// temporary test only
	threatfeed.AddTestIP("::1")
	fmt.Println("Threat feed indicators loaded:", threatfeed.Count())

	// Initialize event pipeline
	fmt.Println("Initializing event pipeline...")
	pipeline.InitEventQueue(EventQueueSize)

	// Start worker pool
	fmt.Println("Starting worker pool...")
	pipeline.StartWorkerPool(WorkerCount, processEvent)

	// Start metrics reporter
	go metrics.StartMetricsReporter()

	// Start HTTP collector
	go collector.StartHTTPServer()

	// Start API server
	go api.StartAPIServer()

	fmt.Println("SentinelX running")

	select {}
}

func processEvent(event models.SecurityEvent) {
	// Save raw event
	storage.SaveEvent(event)

	// Metrics
	metrics.RecordEvent(event.SourceIP, event.EventType)

	// Correlation tracking
	correlation.RecordEvent(event.SourceIP, event.EventType)

	eventTime := time.Unix(0, event.Timestamp).UTC()

	// External threat feed check
	if threatfeed.IsMalicious(event.SourceIP) {
		alert := models.Alert{
			ID:          generateMainAlertID(),
			Timestamp:   eventTime,
			Type:        "threat_intel_match",
			Severity:    models.SeverityCritical,
			SourceIP:    event.SourceIP,
			Description: "Source IP matched external threat intelligence feed",
			ThreatScore: 0.98,
			Status:      models.AlertStatusNew,
			Metadata: map[string]interface{}{
				"matched_ip": event.SourceIP,
				"source":     "external_threat_feed",
			},
		}

		select {
		case detection.AlertQueue <- alert:
		default:
			fmt.Println("Alert queue full — dropping external threat intel alert")
		}
	}

	// Correlation-based multi-stage detection
	if correlation.DetectMultiStage(event.SourceIP) {
		if alert := correlation.BuildMultiStageAlert(event.SourceIP); alert != nil {
			select {
			case detection.AlertQueue <- *alert:
			default:
				fmt.Println("Alert queue full — dropping multi-stage alert")
			}
		}
	}

	// Local threat intel engine
	if threatintel.IsMaliciousIP(event.SourceIP) {
		alert := models.Alert{
			ID:          generateMainAlertID(),
			Timestamp:   eventTime,
			Type:        "threat_intel_match",
			Severity:    models.SeverityCritical,
			SourceIP:    event.SourceIP,
			Description: "Connection from known malicious IP",
			ThreatScore: 0.98,
			Status:      models.AlertStatusNew,
			Metadata: map[string]interface{}{
				"matched_ip": event.SourceIP,
				"source":     "local_threat_intel",
			},
		}

		select {
		case detection.AlertQueue <- alert:
		default:
			fmt.Println("Alert queue full — dropping local threat intel alert")
		}
	}

	// Rule engine processing
	if alert := ruleengine.ProcessEvent(event); alert != nil {
		select {
		case detection.AlertQueue <- *alert:
		default:
			fmt.Println("Alert queue full — dropping rule-engine alert")
		}
	}

	// Port scan detector
	if alert := detection.ScanDetector.ProcessEvent(event); alert != nil {
		select {
		case detection.AlertQueue <- *alert:
		default:
			fmt.Println("Alert queue full — dropping port-scan alert")
		}
	}

	// WAF detector
	if alert := detection.WAF.ProcessEvent(event); alert != nil {
		select {
		case detection.AlertQueue <- *alert:
		default:
			fmt.Println("Alert queue full — dropping WAF alert")
		}
	}

	// Threat intel detector
	if alert := detection.ThreatIntel.ProcessEvent(event); alert != nil {
		select {
		case detection.AlertQueue <- *alert:
		default:
			fmt.Println("Alert queue full — dropping detector threat-intel alert")
		}
	}
}

func generateMainAlertID() string {
	return fmt.Sprintf("ALT-%d", time.Now().UnixNano())
}