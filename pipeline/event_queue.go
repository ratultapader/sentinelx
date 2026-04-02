package pipeline

import (
	"context"
	"fmt"

	// "sentinelx/detection"
	"sentinelx/models"
	"sentinelx/multi_tenant"
)

// Global event queue
var EventQueue chan models.SecurityEvent

// ===============================
// INIT PIPELINE
// ===============================
func StartPipeline(handler func(models.SecurityEvent)) {
	InitEventQueue(1000)
	StartWorkerPool(5, handler)
}

// Initialize queue
func InitEventQueue(size int) {
	EventQueue = make(chan models.SecurityEvent, size)
}

// ===============================
// PUBLISH EVENT
// ===============================
func PublishEvent(ctx context.Context, event models.SecurityEvent) {
	


	tenantID := multi_tenant.TenantIDFromContext(ctx)

	fmt.Println("📥 PublishEvent CALLED >>>", event)

	if tenantID == "" {
		fmt.Println("missing tenant_id, dropping event")
		return
	}

	event.TenantID = tenantID

	select {
	case EventQueue <- event:
	default:
		fmt.Println("Event queue full, dropping event")
	}
}

// ===============================
// WORKER POOL (FINAL FIX)
// ===============================
func StartWorkerPool(workerCount int, handler func(models.SecurityEvent)) {

	for i := 0; i < workerCount; i++ {

		go func(workerID int) {

			for event := range EventQueue {

				fmt.Println("DEBUG: pipeline worker received event")
				fmt.Println("Worker", workerID, "processing event (tenant:", event.TenantID, ")")

// 				// ✅ KEEP YOUR WAF (as you requested)
// 				alert := detection.WAF.ProcessEvent(event)
// if alert != nil {
// 	fmt.Println("🚨 ALERT TRIGGERED:", alert.Type, alert.SourceIP)

// 	// 🔥 THIS IS THE MISSING LINK
// 	ctx := context.WithValue(context.Background(), "tenant_id", event.TenantID)

// // 🔥 pass target from event metadata
// if t, ok := event.Metadata["target"].(string); ok {
// 	ctx = context.WithValue(ctx, "target", t)
// }

// detection.GenerateAlert(
// 	ctx,
// 	alert.Type,
// 	alert.SourceIP,
// 	"Auto generated from pipeline",
// )
// }

				// 🔥 ADD THIS (CRITICAL FIX)
				handler(event)
			}

		}(i)
	}
}