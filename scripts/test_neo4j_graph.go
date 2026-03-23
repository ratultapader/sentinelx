package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"sentinelx/storage"
)

func main() {
	ctx := context.Background()

	client, err := storage.NewNeo4jClient("neo4j://localhost:7687", "neo4j", "password")
	if err != nil {
		log.Fatalf("neo4j client create failed: %v", err)
	}
	defer client.Close(ctx)

	if err := client.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("neo4j connectivity failed: %v", err)
	}

	store := storage.NewNeo4jGraphStore(client, "neo4j")

	if err := store.EnsureConstraints(ctx); err != nil {
		log.Fatalf("ensure constraints failed: %v", err)
	}

	ingestor := storage.NewNeo4jIngestor(store)

	record := storage.AttackGraphRecord{
		AlertID:        "alert_1001",
		Timestamp:      time.Now().UTC(),
		SourceIP:       "192.168.1.5",
		Server:         "payments-api",
		APIEndpoint:    "/login",
		EventType:      "sql_injection",
		Severity:       "critical",
		ThreatScore:    0.94,
		ResponseAction: "ip_block",
	}

	if err := ingestor.IngestAttackRecord(ctx, record); err != nil {
		log.Fatalf("graph ingest failed: %v", err)
	}

	paths, err := store.GetAttackPathByIP(ctx, "192.168.1.5")
	if err != nil {
		log.Fatalf("graph query failed: %v", err)
	}

	fmt.Println("Neo4j graph test successful")
	fmt.Println("Paths found:", len(paths))
	for i, p := range paths {
		fmt.Printf("Path %d: %+v\n", i+1, p)
	}
}