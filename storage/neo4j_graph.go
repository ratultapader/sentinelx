package storage

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jGraphStore struct {
	client   *Neo4jClient
	database string
}

func NewNeo4jGraphStore(client *Neo4jClient, database string) *Neo4jGraphStore {
	return &Neo4jGraphStore{
		client:   client,
		database: database,
	}
}

func (s *Neo4jGraphStore) EnsureConstraints(ctx context.Context) error {
	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	queries := []string{
		"CREATE CONSTRAINT attacker_ip_key IF NOT EXISTS FOR (n:AttackerIP) REQUIRE n.key IS UNIQUE",
		"CREATE CONSTRAINT server_key IF NOT EXISTS FOR (n:Server) REQUIRE n.key IS UNIQUE",
		"CREATE CONSTRAINT container_key IF NOT EXISTS FOR (n:Container) REQUIRE n.key IS UNIQUE",
		"CREATE CONSTRAINT api_endpoint_key IF NOT EXISTS FOR (n:APIEndpoint) REQUIRE n.key IS UNIQUE",
		"CREATE CONSTRAINT alert_key IF NOT EXISTS FOR (n:Alert) REQUIRE n.key IS UNIQUE",
		"CREATE CONSTRAINT response_action_key IF NOT EXISTS FOR (n:ResponseAction) REQUIRE n.key IS UNIQUE",
	}

	for _, q := range queries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
			_, err := tx.Run(ctx, q, nil)
			return nil, err
		})
		if err != nil {
			return fmt.Errorf("ensure neo4j constraint failed: %w", err)
		}
	}

	return nil
}

func (s *Neo4jGraphStore) UpsertNode(ctx context.Context, node GraphNode) error {
	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MERGE (n:%s {key: $key})
		SET n += $props
	`, node.Label)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"key":   node.Key,
			"props": node.Properties,
		})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("upsert node %s/%s: %w", node.Label, node.Key, err)
	}

	return nil
}

func (s *Neo4jGraphStore) UpsertRelationship(ctx context.Context, rel GraphRelationship) error {
	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MATCH (a:%s {key: $fromKey})
		MATCH (b:%s {key: $toKey})
		MERGE (a)-[r:%s]->(b)
		SET r += $props
	`, rel.FromLabel, rel.ToLabel, rel.Type)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"fromKey": rel.FromKey,
			"toKey":   rel.ToKey,
			"props":   rel.Properties,
		})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf(
			"upsert relationship %s %s/%s -> %s/%s: %w",
			rel.Type, rel.FromLabel, rel.FromKey, rel.ToLabel, rel.ToKey, err,
		)
	}

	return nil
}

func (s *Neo4jGraphStore) GetAttackPathByIP(ctx context.Context, sourceIP string) ([]map[string]interface{}, error) {
	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := `
		MATCH p=(a:AttackerIP {key: $sourceIP})-[*1..4]->(n)
		RETURN p
		LIMIT 25
	`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		res, err := tx.Run(ctx, query, map[string]interface{}{
			"sourceIP": sourceIP,
		})
		if err != nil {
			return nil, err
		}

		paths := make([]map[string]interface{}, 0)
		for res.Next(ctx) {
			record := res.Record()
			paths = append(paths, map[string]interface{}{
				"path": record.Values[0],
			})
		}
		return paths, res.Err()
	})
	if err != nil {
		return nil, fmt.Errorf("get attack path by ip %s: %w", sourceIP, err)
	}

	if result == nil {
		return []map[string]interface{}{}, nil
	}
	return result.([]map[string]interface{}), nil
}

func (s *Neo4jGraphStore) ExportGraphBySourceIP(ctx context.Context, sourceIP string) (GraphView, error) {
	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := `
		MATCH (src:AttackerIP {key: $sourceIP})-[*0..3]->(m)
		WITH collect(DISTINCT src) + collect(DISTINCT m) AS reachable
		UNWIND reachable AS a
		OPTIONAL MATCH (a)-[r]->(b)
		WHERE b IN reachable
		RETURN DISTINCT a, r, b
	`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		res, err := tx.Run(ctx, query, map[string]interface{}{
			"sourceIP": sourceIP,
		})
		if err != nil {
			return nil, err
		}

		graph := GraphView{
			Nodes: []GraphNodeView{},
			Links: []GraphLinkView{},
		}

		nodeSeen := map[string]struct{}{}
		linkSeen := map[string]struct{}{}

		for res.Next(ctx) {
			record := res.Record()

			aRaw, _ := record.Get("a")
			bRaw, _ := record.Get("b")
			rRaw, _ := record.Get("r")

			if aNode, ok := aRaw.(neo4j.Node); ok {
				addNodeView(seenNodeArg{}, nodeSeen, &graph, aNode)
			}
			if bNode, ok := bRaw.(neo4j.Node); ok {
				addNodeView(seenNodeArg{}, nodeSeen, &graph, bNode)
			}

			if rRaw != nil {
				if rel, ok := rRaw.(neo4j.Relationship); ok {
					aNode, aOK := aRaw.(neo4j.Node)
					bNode, bOK := bRaw.(neo4j.Node)
					if aOK && bOK {
						aKey, _ := aNode.Props["key"].(string)
						bKey, _ := bNode.Props["key"].(string)
						if aKey != "" && bKey != "" {
							linkKey := aKey + "|" + rel.Type + "|" + bKey
							if _, exists := linkSeen[linkKey]; !exists {
								linkSeen[linkKey] = struct{}{}
								graph.Links = append(graph.Links, GraphLinkView{
									Source:     aKey,
									Target:     bKey,
									Type:       rel.Type,
									Properties: rel.Props,
								})
							}
						}
					}
				}
			}
		}

		return graph, res.Err()
	})
	if err != nil {
		return GraphView{}, fmt.Errorf("export graph by source ip %s: %w", sourceIP, err)
	}

	if result == nil {
		return GraphView{}, nil
	}
	return result.(GraphView), nil
}

type seenNodeArg struct{}

func addNodeView(_ seenNodeArg, seen map[string]struct{}, graph *GraphView, node neo4j.Node) {
	key, _ := node.Props["key"].(string)
	if key == "" {
		return
	}
	if _, exists := seen[key]; exists {
		return
	}
	seen[key] = struct{}{}

	label := ""
	if len(node.Labels) > 0 {
		label = node.Labels[0]
	}

	graph.Nodes = append(graph.Nodes, GraphNodeView{
		Label:      label,
		Key:        key,
		Properties: node.Props,
	})
}