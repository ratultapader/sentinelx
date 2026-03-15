package main

import (
	"time"

	"sentinelx/collector"
	"sentinelx/detection"
)

func main() {

	detection.InitAlertEngine(100)

	go detection.StartAlertProcessor()

	event := collector.NewSecurityEvent("http_request")

	event.SourceIP = "192.168.1.5"

	event.Metadata["path"] = "/search?q=<script>alert(1)</script>"

	detection.WAF.ProcessEvent(event)

	time.Sleep(1 * time.Second)

}