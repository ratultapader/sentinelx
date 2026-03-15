package pipeline

import (
	"fmt"
	"sentinelx/models"
)

// Global event queue
var EventQueue chan models.SecurityEvent

// Initialize queue
func InitEventQueue(size int) {
	EventQueue = make(chan models.SecurityEvent, size)
}

// Publish event into queue
func PublishEvent(event models.SecurityEvent) {

	select {

	case EventQueue <- event:

	default:
		// queue full → drop event
		fmt.Println("Event queue full, dropping event")

	}
}

// Worker pool implementation
func StartWorkerPool(workerCount int, handler func(models.SecurityEvent)) {

	for i := 0; i < workerCount; i++ {

		go func(workerID int) {

			for event := range EventQueue {

				fmt.Println("Worker", workerID, "processing event")

				handler(event)

			}

		}(i)

	}
}