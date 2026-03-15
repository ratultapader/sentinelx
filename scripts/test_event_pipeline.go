package main

import (
	"fmt"
	"time"

	"sentinelx/collector"
	"sentinelx/pipeline"
)

func main() {

	pipeline.InitEventQueue(1000)

	go pipeline.StartEventConsumer(processEvent)

	for i := 0; i < 10; i++ {

		event := collector.NewSecurityEvent("test_event")

		event.SourceIP = "192.168.1.10"

		pipeline.PublishEvent(event)

	}

	time.Sleep(2 * time.Second)

}

func processEvent(event collector.SecurityEvent) {

	fmt.Println("Processing event:", event.EventType, event.SourceIP)

}