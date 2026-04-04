package storage

import (
	"context"
	"log"
)

var GraphIngestor *Neo4jIngestor
var GraphStore *Neo4jGraphStore
var graphClient *Neo4jClient

func InitNeo4jGraph(uri, username, password, database string) {
	client, err := NewNeo4jClient()
	if err != nil {
		log.Printf("WARNING: Neo4j init failed: %v", err)
		return
	}

	if err := client.VerifyConnectivity(context.Background()); err != nil {
		log.Printf("WARNING: Neo4j connectivity failed: %v", err)
		_ = client.Close(context.Background())
		return
	}

	store := NewNeo4jGraphStore(client, database)
	if err := store.EnsureConstraints(context.Background()); err != nil {
		log.Printf("WARNING: Neo4j constraints failed: %v", err)
		_ = client.Close(context.Background())
		return
	}

	graphClient = client
	GraphStore = store
	GraphIngestor = NewNeo4jIngestor(store)
	log.Println("Neo4j Graph Engine initialized")
}

func CloseNeo4jGraph() {
	if graphClient != nil {
		_ = graphClient.Close(context.Background())
	}
}