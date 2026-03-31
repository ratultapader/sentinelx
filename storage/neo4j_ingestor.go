package storage

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Neo4jIngestor struct {
	store *Neo4jGraphStore
}

func getResponseAction(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "critical":
		return "block_ip"
	case "high":
		return "rate_limit"
	case "medium":
		return "monitor"
	default:
		return "alert_only"
	}
}

func NewNeo4jIngestor(store *Neo4jGraphStore) *Neo4jIngestor {
	return &Neo4jIngestor{
		store: store,
	}
}

func (i *Neo4jIngestor) IngestAttackRecord(ctx context.Context, record AttackGraphRecord) error {
	tenantID := strings.TrimSpace(record.TenantID)
	if tenantID == "" {
		return fmt.Errorf("tenant_id is required")
	}

	sourceIPKey := strings.TrimSpace(record.SourceIP)
	alertKey := strings.TrimSpace(record.AlertID)

	if sourceIPKey == "" || alertKey == "" {
		return fmt.Errorf("source ip and alert id are required")
	}

	attackerNode := GraphNode{
		Label: NodeAttackerIP,
		Key:   sourceIPKey,
		Properties: map[string]interface{}{
			"key":       sourceIPKey,
			"ip":        sourceIPKey,
			"tenant_id": tenantID,
			"last_seen": record.Timestamp.Format(time.RFC3339Nano),
		},
	}
	if err := i.store.UpsertNode(ctx, attackerNode); err != nil {
		return err
	}

	alertNode := GraphNode{
		Label: NodeAlert,
		Key:   alertKey,
		Properties: map[string]interface{}{
			"key":                alertKey,
			"alert_id":           record.AlertID,
			"tenant_id":          tenantID,
			"event_type":         record.EventType,
			"severity":           record.Severity,
			"threat_score":       record.ThreatScore,
			"timestamp":          record.Timestamp.Format(time.RFC3339Nano),
			"mitre_tactic":       record.MitreTactic,
			"mitre_technique":    record.MitreTechnique,
			"mitre_technique_id": record.MitreTechniqueID,
		},
	}
	if err := i.store.UpsertNode(ctx, alertNode); err != nil {
		return err
	}

	if err := i.store.UpsertRelationship(ctx, GraphRelationship{
		Type:      RelTriggered,
		FromLabel: NodeAttackerIP,
		FromKey:   sourceIPKey,
		ToLabel:   NodeAlert,
		ToKey:     alertKey,
		Properties: map[string]interface{}{
			"timestamp": record.Timestamp.Format(time.RFC3339Nano),
		},
	}); err != nil {
		return err
	}

	if strings.TrimSpace(record.Server) != "" {
		serverKey := strings.TrimSpace(record.Server)
		serverNode := GraphNode{
			Label: NodeServer,
			Key:   serverKey,
			Properties: map[string]interface{}{
				"key":       serverKey,
				"name":      serverKey,
				"tenant_id": tenantID,
			},
		}
		if err := i.store.UpsertNode(ctx, serverNode); err != nil {
			return err
		}

		if err := i.store.UpsertRelationship(ctx, GraphRelationship{
			Type:      RelAttacked,
			FromLabel: NodeAttackerIP,
			FromKey:   sourceIPKey,
			ToLabel:   NodeServer,
			ToKey:     serverKey,
			Properties: map[string]interface{}{
				"event_type":   record.EventType,
				"severity":     record.Severity,
				"threat_score": record.ThreatScore,
				"timestamp":    record.Timestamp.Format(time.RFC3339Nano),
			},
		}); err != nil {
			return err
		}
	}

	if strings.TrimSpace(record.APIEndpoint) != "" {
		endpointKey := strings.TrimSpace(record.APIEndpoint)
		endpointNode := GraphNode{
			Label: NodeAPIEndpoint,
			Key:   endpointKey,
			Properties: map[string]interface{}{
				"key":       endpointKey,
				"endpoint":  endpointKey,
				"tenant_id": tenantID,
			},
		}
		if err := i.store.UpsertNode(ctx, endpointNode); err != nil {
			return err
		}

		if strings.TrimSpace(record.Server) != "" {
			serverKey := strings.TrimSpace(record.Server)
			if serverKey != endpointKey {
				if err := i.store.UpsertRelationship(ctx, GraphRelationship{
					Type:      RelConnectedTo,
					FromLabel: NodeServer,
					FromKey:   serverKey,
					ToLabel:   NodeAPIEndpoint,
					ToKey:     endpointKey,
					Properties: map[string]interface{}{
						"timestamp": record.Timestamp.Format(time.RFC3339Nano),
					},
				}); err != nil {
					return err
				}
			}
		}

		if err := i.store.UpsertRelationship(ctx, GraphRelationship{
			Type:      RelTargeted,
			FromLabel: NodeAlert,
			FromKey:   alertKey,
			ToLabel:   NodeAPIEndpoint,
			ToKey:     endpointKey,
			Properties: map[string]interface{}{
				"event_type": record.EventType,
			},
		}); err != nil {
			return err
		}
	}

	if strings.TrimSpace(record.Container) != "" {
		containerKey := strings.TrimSpace(record.Container)
		containerNode := GraphNode{
			Label: NodeContainer,
			Key:   containerKey,
			Properties: map[string]interface{}{
				"key":       containerKey,
				"name":      containerKey,
				"tenant_id": tenantID,
			},
		}
		if err := i.store.UpsertNode(ctx, containerNode); err != nil {
			return err
		}

		if strings.TrimSpace(record.Server) != "" {
			serverKey := strings.TrimSpace(record.Server)
			if serverKey != containerKey {
				if err := i.store.UpsertRelationship(ctx, GraphRelationship{
					Type:      RelConnectedTo,
					FromLabel: NodeServer,
					FromKey:   serverKey,
					ToLabel:   NodeContainer,
					ToKey:     containerKey,
					Properties: map[string]interface{}{
						"timestamp": record.Timestamp.Format(time.RFC3339Nano),
					},
				}); err != nil {
					return err
				}
			}
		}

		if record.EventType == "container_escape" || record.EventType == "runtime_compromise" {
			if err := i.store.UpsertRelationship(ctx, GraphRelationship{
				Type:      RelExploited,
				FromLabel: NodeAttackerIP,
				FromKey:   sourceIPKey,
				ToLabel:   NodeContainer,
				ToKey:     containerKey,
				Properties: map[string]interface{}{
					"event_type":   record.EventType,
					"severity":     record.Severity,
					"threat_score": record.ThreatScore,
					"timestamp":    record.Timestamp.Format(time.RFC3339Nano),
				},
			}); err != nil {
				return err
			}
		}
	}

	actionType := getResponseAction(record.Severity)
	responseKey := fmt.Sprintf("%s:%s", alertKey, actionType)

	responseNode := GraphNode{
		Label: NodeResponseAction,
		Key:   responseKey,
		Properties: map[string]interface{}{
			"key":         responseKey,
			"action_type": actionType,
			"tenant_id":   tenantID,
			"timestamp":   record.Timestamp.Format(time.RFC3339Nano),
		},
	}
	if err := i.store.UpsertNode(ctx, responseNode); err != nil {
		return err
	}

	if err := i.store.UpsertRelationship(ctx, GraphRelationship{
		Type:      RelMitigatedBy,
		FromLabel: NodeAlert,
		FromKey:   alertKey,
		ToLabel:   NodeResponseAction,
		ToKey:     responseKey,
		Properties: map[string]interface{}{
			"timestamp": record.Timestamp.Format(time.RFC3339Nano),
		},
	}); err != nil {
		return err
	}

	return nil
}
