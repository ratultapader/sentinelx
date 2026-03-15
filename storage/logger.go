package storage

import (
	"encoding/json"
	"os"
	"sync"

	"sentinelx/models"
)

var logFile *os.File
var mutex sync.Mutex

// InitLogger initializes the security log file
func InitLogger(path string) error {

	var err error

	logFile, err = os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	return err
}

// LogAlert writes alerts to log file
func LogAlert(alert models.Alert) {

	mutex.Lock()
	defer mutex.Unlock()

	jsonData, err := json.Marshal(alert)
	if err != nil {
		return
	}

	logFile.Write(jsonData)
	logFile.Write([]byte("\n"))
}
