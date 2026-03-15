package main

import (
	"sentinelx/detection"
)

func main() {

	detection.InitAlertEngine(100)

	go detection.StartAlertProcessor()

	detection.GenerateAlert(
		"HIGH",
		"sql_injection",
		"192.168.1.5",
		"SQL injection attempt detected",
	)

	select {}

}