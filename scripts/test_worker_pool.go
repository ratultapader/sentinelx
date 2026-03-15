package main

import (
	"fmt"
	"sentinelx/models"
	"sentinelx/pipeline"
)

func main() {

	pipeline.InitEventQueue(10000)

	pipeline.StartWorkerPool(5, processEvent)

	for i := 0; i < 1000; i++ {

		event := models.NewSecurityEvent("test_event")

		pipeline.PublishEvent(event)

	}

	select {}

}

func processEvent(event models.SecurityEvent) {

	fmt.Println("Processing event:", event.EventType)

}
