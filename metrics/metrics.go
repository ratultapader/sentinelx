package metrics

import "sync"

var mutex sync.Mutex

var TotalEvents int
var TotalAlerts int

var AttackTypes = make(map[string]int)
var AttackerIPs = make(map[string]int)

func RecordEvent(eventType string, ip string) {

	mutex.Lock()
	defer mutex.Unlock()

	TotalEvents++

	AttackTypes[eventType]++
	AttackerIPs[ip]++
}

func RecordAlert() {

	mutex.Lock()
	defer mutex.Unlock()

	TotalAlerts++
}