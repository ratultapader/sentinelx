package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jClient struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jClient() (*Neo4jClient, error) {

	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	if uri == "" {
		return nil, fmt.Errorf("NEO4J_URI not set")
	}

	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("create neo4j driver: %w", err)
	}

	return &Neo4jClient{
		driver: driver,
	}, nil
}

func (n *Neo4jClient) VerifyConnectivity(ctx context.Context) error {
	if err := n.driver.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("verify neo4j connectivity: %w", err)
	}
	return nil
}

func (n *Neo4jClient) Close(ctx context.Context) error {
	return n.driver.Close(ctx)
}

func (n *Neo4jClient) Driver() neo4j.DriverWithContext {
	return n.driver
}