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

	sub := storage.RDB.Subscribe(ctx, "alerts_channel")

	ch := sub.Channel()

	log.Println("📡 Redis subscriber started...")

	for msg := range ch {
		var alert models.Alert

		err := json.Unmarshal([]byte(msg.Payload), &alert)
		if err != nil {
			continue
		}

		log.Println("⚡ New alert from Redis:", alert.SourceIP)

		// 🔥 SEND TO WEBSOCKET
		BroadcastAlert(alert)
	}
}