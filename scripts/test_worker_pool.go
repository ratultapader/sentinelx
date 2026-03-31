package main

import (
	"context"
	"fmt"

	"sentinelx/detection"
	"sentinelx/models"
	"sentinelx/multi_tenant"
	"sentinelx/pipeline"
)

func main() {
	pipeline.InitEventQueue(10000)

	pipeline.StartWorkerPool(5, func(event models.SecurityEvent) {
		alert := detection.WAF.ProcessEvent(event)
		if alert != nil {
			fmt.Println("ALERT TRIGGERED:", alert.Type, alert.SourceIP)
		}
	})

	for i := 0; i < 1000; i++ {
		event := models.NewSecurityEvent("test_event")
		ctx := multi_tenant.WithTenantID(context.Background(), "t1")
		pipeline.PublishEvent(ctx, event)
	}

	select {}
}
