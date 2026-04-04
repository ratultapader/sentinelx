package configs

import (
	"encoding/json"
)

func Log(level string, message string, fields map[string]interface{}) {
	logData := map[string]interface{}{
		"level":   level,
		"message": message,
	}

	for k, v := range fields {
		logData[k] = v
	}

	jsonLog, _ := json.Marshal(logData)
	Logger.Println(string(jsonLog))
}