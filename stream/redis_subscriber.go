package stream

import (
	"context"
	"encoding/json"
	"log"

	"sentinelx/models"
	"sentinelx/storage"
)

func StartRedisSubscriber() {
	ctx := context.Background()

	// 🔥 FIX 1: Check Redis
	if storage.RDB == nil {
		log.Println("⚠️ Redis not available — subscriber disabled")
		return
	}

	// 🔥 FIX 2: Safe subscribe
	sub := storage.RDB.Subscribe(ctx, "alerts_channel")

	ch := sub.Channel()

	log.Println("📡 Redis subscriber started...")

	for msg := range ch {
		var alert models.Alert

		err := json.Unmarshal([]byte(msg.Payload), &alert)
		if err != nil {
			log.Println("⚠️ Failed to parse alert:", err)
			continue
		}

		log.Println("⚡ New alert from Redis:", alert.SourceIP)

		// 🔥 FIX 3: Safe WebSocket call (optional but good)
		BroadcastAlert(alert)
	}
}