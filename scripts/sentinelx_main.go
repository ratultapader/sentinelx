package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"sentinelx/api"
	"sentinelx/app"
	"sentinelx/stream"
	// "sentinelx/collector"
	"sentinelx/configs"
	"sentinelx/correlation"
	"sentinelx/detection"
	"sentinelx/metrics"
	"sentinelx/models"
	"sentinelx/pipeline"
	"sentinelx/repository"
	"sentinelx/response"
	"sentinelx/ruleengine"
	"sentinelx/service"
	"sentinelx/storage"
	"sentinelx/threatfeed"
	"sentinelx/threatintel"
	// "sentinelx/configs"
)

const (
	EventQueueSize = 10000
	AlertQueueSize = 1000
	WorkerCount    = 5
)

var alertCache = make(map[string]*models.Alert)
var alertCacheMutex sync.Mutex

func main() {
	fmt.Println("Starting SentinelX Security Platform")
	configs.InitLogger()
	configs.InitMetrics()   // ✅ ADD THIS

	// ✅ BOOTSTRAP (FINAL FIX)
	deps := app.Bootstrap()
	_ = deps

	cfg := configs.Load()

	// ===============================
	// THREAT INTEL
	// ===============================
	err := threatintel.LoadThreatFeed("data/malicious_ips.txt")
	if err != nil {
		panic(err)
	}

	// ===============================
	// LOGGER + DB
	// ===============================
	fmt.Println("Initializing event logger...")
	err = storage.InitLogger("logs/security_events.json")
	if err != nil {
		panic(err)
	}

	err = storage.InitDB()
	if err != nil {
		panic(err)
	}

storage.InitRedis()

if storage.RDB != nil {
    go stream.StartRedisSubscriber()
}

	// ===============================
	// ELASTICSEARCH (CONFIG-DRIVEN)
	// ===============================
	err = storage.InitElasticsearch(storage.ElasticsearchConfig{
		Enabled:   true,
		Addresses: []string{cfg.ElasticsearchURL},
		Username:  "",
		Password:  "",
	})
	if err != nil {
		fmt.Println("WARNING: Elasticsearch init failed:", err)
	} else {
		fmt.Println("Elasticsearch forensic store initialized")
	}

	// ===============================
	// NEO4J (CONFIG-DRIVEN)
	// ===============================
	storage.InitNeo4jGraph(
		cfg.Neo4jURI,
		cfg.Neo4jUsername,
		cfg.Neo4jPassword,
		"neo4j",
	)
	defer storage.CloseNeo4jGraph()

	// ===============================
	// RULE ENGINE
	// ===============================
	fmt.Println("Loading detection rules...")
	err = ruleengine.LoadRules()
	if err != nil {
		panic(err)
	}

	// ===============================
	// ALERT ENGINE
	// ===============================
	fmt.Println("Initializing alert engine...")
	detection.InitAlertEngine(AlertQueueSize)
	go detection.StartAlertProcessor()

	// ===============================
	// RESPONSE ENGINE
	// ===============================
	fmt.Println("Initializing response engine...")
	response.InitResponseEngine(1000)
	go response.StartActionProcessor()

	fmt.Println("Initializing firewall executor...")
	response.InitFirewallExecutor(1000)
	blocker := response.NewFirewallBlocker(true)
	response.StartFirewallExecutor(blocker)

	fmt.Println("Initializing rate limit executor...")
	response.InitRateLimitExecutor(1000)
	limiter := response.NewRateLimiter(true)
	response.StartRateLimitExecutor(limiter, 20, 40)

	fmt.Println("Initializing kubernetes executor...")
	response.InitKubernetesExecutor(1000)
	k8sController := response.NewKubernetesController(true)
	response.StartKubernetesExecutor(k8sController)

	// ===============================
	// THREAT FEED
	// ===============================
	fmt.Println("Starting threat feed updater...")
	threatfeed.StartThreatFeedUpdater()

	threatfeed.AddTestIP("::1")
	fmt.Println("Threat feed indicators loaded:", threatfeed.Count())

	// 🔥 CLEAN ARCHITECTURE (NEW)
	alertRepo := repository.NewAlertRepository()
	alertService := service.NewAlertService(alertRepo)

	// ===============================
	// PIPELINE
	// ===============================
	fmt.Println("Initializing event pipeline...")
	pipeline.InitEventQueue(EventQueueSize)

	fmt.Println("Starting worker pool...")
	pipeline.StartWorkerPool(WorkerCount, func(event models.SecurityEvent) {
		processEvent(event, alertService)
	})

	// ===============================
	// METRICS
	// ===============================
	go metrics.StartMetricsReporter()

	// ===============================
	// SERVERS
	// ===============================
	// go collector.StartHTTPServer()
	go api.StartAPIServer(alertService)

	configs.Log("INFO", "SentinelX started", map[string]interface{}{})

	select {}
}

// ===============================
// EVENT PROCESSOR (UNCHANGED)
// ===============================

func processEvent(event models.SecurityEvent, alertService *service.AlertService) {

	ctx := context.Background()

	fmt.Println("DEBUG: processEvent CALLED")

	storage.SaveEvent(event)
	correlation.RecordEvent(event.SourceIP, event.EventType)

	eventTime := time.Unix(0, event.Timestamp).UTC()

	if threatfeed.IsMalicious(event.SourceIP) {
		alert := models.Alert{
			ID:          generateMainAlertID(),
			TenantID:    event.TenantID,
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
			metrics.RecordAlert(alert.Type)
			alertService.ProcessAlert(ctx, alert)
		default:
			fmt.Println("Alert queue full — dropping alert")
		}
	}

	if correlation.DetectMultiStage(event.SourceIP) {
		if alert := correlation.BuildMultiStageAlert(event.SourceIP, event.TenantID); alert != nil {
			select {
			case detection.AlertQueue <- *alert:
				metrics.RecordAlert(alert.Type)
				alertService.ProcessAlert(ctx, *alert)
			default:
				fmt.Println("Alert queue full — dropping alert")
			}
		}
	}

	if threatintel.IsMaliciousIP(event.SourceIP) {
		alert := models.Alert{
			ID:          generateMainAlertID(),
			TenantID:    event.TenantID,
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
			metrics.RecordAlert(alert.Type)
			alertService.ProcessAlert(ctx, alert)
		default:
			fmt.Println("Alert queue full — dropping alert")
		}
	}

	ruleAlert := ruleengine.ProcessEvent(event)
	if ruleAlert != nil {
		select {
		case detection.AlertQueue <- *ruleAlert:
			metrics.RecordAlert(ruleAlert.Type)
			alertService.ProcessAlert(ctx, *ruleAlert)
		default:
			fmt.Println("Alert queue full — dropping alert")
		}
	}

	detected := event.EventType
	if ruleAlert != nil {
		if v, ok := ruleAlert.Metadata["detected_type"].(string); ok {
			detected = v
		}
	}

	metrics.RecordEvent(event.SourceIP, event.EventType, detected)

	if alert := detection.ScanDetector.ProcessEvent(event); alert != nil {

		// skip SQLi (handled by WAF)
		if alert.Type == "sql_injection" {
			return
		}

		select {
		case detection.AlertQueue <- *alert:
			metrics.RecordAlert(alert.Type)
			alertService.ProcessAlert(ctx, *alert)
		default:
			fmt.Println("Alert queue full � dropping alert")
		}
	}

	if alert := detection.WAF.ProcessEvent(event); alert != nil {

		// 🔥 PREVENT DUPLICATE PROCESSING
		if event.EventType == "sql_injection" {
			return
		}

		key := event.SourceIP + "_" + alert.Type + "_" + fmt.Sprint(event.Timestamp)

		// if already exists -> increase count
		if existing, ok := alertCache[key]; ok {

			if existing.Metadata == nil {
				existing.Metadata = make(map[string]interface{})
			}

			count, _ := existing.Metadata["count"].(int)
			existing.Metadata["count"] = count + 1

			return
		}

		// first time alert
		alert.Metadata["count"] = 1
		alertCache[key] = alert

		select {
		case detection.AlertQueue <- *alert:
			metrics.RecordAlert(alert.Type)
			alertService.ProcessAlert(ctx, *alert)
		default:
			fmt.Println("Alert queue full — dropping alert")
		}
	}

	if alert := detection.ThreatIntel.ProcessEvent(event); alert != nil {
		select {
		case detection.AlertQueue <- *alert:
			metrics.RecordAlert(alert.Type)
			alertService.ProcessAlert(ctx, *alert)
		default:
			fmt.Println("Alert queue full — dropping alert")
		}
	}
}

func generateMainAlertID() string {
	return fmt.Sprintf("ALT-%d-%d", time.Now().UnixNano(), rand.Intn(10000))
}
