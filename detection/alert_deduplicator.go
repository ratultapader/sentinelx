package detection

import (
	"fmt"
	"sync"
	"time"

	"sentinelx/models"
)

type alertRecord struct {
	lastSeen time.Time
	count    int
}

var alertCache = make(map[string]*alertRecord)
var cacheMutex sync.Mutex

const suppressionWindow = 30 * time.Second

func shouldSuppress(alert models.Alert) bool {

	fingerprint := fmt.Sprintf("%s:%s", alert.Type, alert.SourceIP)

	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	record, exists := alertCache[fingerprint]

	if !exists {

		alertCache[fingerprint] = &alertRecord{
			lastSeen: time.Now(),
			count:    1,
		}

		return false
	}

	if time.Since(record.lastSeen) < suppressionWindow {

		record.count++
		return true
	}

	record.lastSeen = time.Now()
	record.count++

	return false
}