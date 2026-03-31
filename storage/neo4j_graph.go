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

//
// ===============================
// UPSERT NODE
// ===============================
//
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

//
// ===============================
// UPSERT RELATIONSHIP (WITH FIX)
// ===============================
//
func (s *Neo4jGraphStore) UpsertRelationship(ctx context.Context, rel GraphRelationship) error {

	// 🔥 FIX: remove self-loop
	if rel.FromKey == rel.ToKey {
		return nil
	}

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

//
// ===============================
// GET ATTACK PATH (CLEAN)
// ===============================
//
func (s *Neo4jGraphStore) GetAttackPathByIP(ctx context.Context, sourceIP string) ([]map[string]interface{}, error) {

	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := `
		MATCH p=(a:AttackerIP {key: $sourceIP})-[*1..4]->(n)
		RETURN p
		LIMIT 10
	`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {

		res, err := tx.Run(ctx, query, map[string]interface{}{
			"sourceIP": sourceIP,
		})
		if err != nil {
			return nil, err
		}

		paths := []map[string]interface{}{}

		for res.Next(ctx) {
			record := res.Record()

			if pathVal, ok := record.Get("p"); ok {
				paths = append(paths, map[string]interface{}{
					"path": pathVal,
				})
			}
		}

		return paths, res.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("get attack path for %s: %w", sourceIP, err)
	}

	if result == nil {
		return []map[string]interface{}{}, nil
	}

	return result.([]map[string]interface{}), nil
}

//
// ===============================
// EXPORT GRAPH (PRO CLEAN)
// ===============================
//
func (s *Neo4jGraphStore) ExportGraphBySourceIP(
	ctx context.Context,
	sourceIP string,
	tenantID string,
) (GraphView, error) {

	session := s.client.Driver().NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: s.database,
	})
	defer session.Close(ctx)

	query := `
		MATCH (src:AttackerIP {key: $sourceIP, tenant_id: $tenantID})
	OPTIONAL MATCH (src)-[r*1..3]->(n)
	WITH collect(DISTINCT src) + collect(DISTINCT n) AS nodes
	UNWIND nodes AS node
	OPTIONAL MATCH (node)-[rel]->(other)
	WHERE other IN nodes
	AND coalesce(node.tenant_id, $tenantID) = $tenantID
	AND coalesce(other.tenant_id, $tenantID) = $tenantID
	RETURN node, rel, other
	`

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {

		res, err := tx.Run(ctx, query, map[string]interface{}{
			"sourceIP": sourceIP,
			"tenantID": tenantID,
		})
		if err != nil {
			return nil, err
		}

		graph := GraphView{
			Nodes: []GraphNodeView{},
			Links: []GraphLinkView{},
		}

		nodeSeen := map[string]bool{}
		linkSeen := map[string]bool{}

		for res.Next(ctx) {
			record := res.Record()

			nRaw, _ := record.Get("node")
			oRaw, _ := record.Get("other")
			rRaw, _ := record.Get("rel")

			// ----- NODE -----
			if n, ok := nRaw.(neo4j.Node); ok {
				key, _ := n.Props["key"].(string)
				if key != "" && !nodeSeen[key] {
					nodeSeen[key] = true

					label := ""
					if len(n.Labels) > 0 {
						label = n.Labels[0]
					}

					graph.Nodes = append(graph.Nodes, GraphNodeView{
						Label:      label,
						Key:        key,
						Properties: n.Props,
					})
				}
			}

			// ----- TARGET NODE -----
			if o, ok := oRaw.(neo4j.Node); ok {
				key, _ := o.Props["key"].(string)
				if key != "" && !nodeSeen[key] {
					nodeSeen[key] = true

					label := ""
					if len(o.Labels) > 0 {
						label = o.Labels[0]
					}

					graph.Nodes = append(graph.Nodes, GraphNodeView{
						Label:      label,
						Key:        key,
						Properties: o.Props,
					})
				}
			}

			// ----- RELATION -----
			if r, ok := rRaw.(neo4j.Relationship); ok {

				fromNode, _ := nRaw.(neo4j.Node)
				toNode, _ := oRaw.(neo4j.Node)

				fromKey, _ := fromNode.Props["key"].(string)
				toKey, _ := toNode.Props["key"].(string)

				if fromKey != "" && toKey != "" {

					linkKey := fromKey + "|" + r.Type + "|" + toKey

					if !linkSeen[linkKey] {
						linkSeen[linkKey] = true

						graph.Links = append(graph.Links, GraphLinkView{
							Source:     fromKey,
							Target:     toKey,
							Type:       r.Type,
							Properties: r.Props,
						})
					}
				}
			}
		}

		return graph, res.Err()
	})

	if err != nil {
		return GraphView{}, fmt.Errorf("export graph error: %w", err)
	}

	if result == nil {
		return GraphView{}, nil
	}

	return result.(GraphView), nil
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
			return fmt.Errorf("constraint failed: %w", err)
		}
	}

	return nil
}