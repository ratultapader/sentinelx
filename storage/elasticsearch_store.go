package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var ESStore *ElasticsearchStore

type ElasticsearchStore struct {
	client *ElasticsearchClient
}

func NewElasticsearchStore(client *ElasticsearchClient) *ElasticsearchStore {
	return &ElasticsearchStore{
		client: client,
	}
}

func InitElasticsearch(cfg ElasticsearchConfig) error {
	if !cfg.Enabled {
		fmt.Println("DEBUG ES: disabled")
		return nil
	}

	fmt.Println("DEBUG ES: creating client")
	client, err := NewElasticsearchClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("DEBUG ES: before ping")
	if err := client.Ping(ctx); err != nil {
		return err
	}
	fmt.Println("DEBUG ES: after ping")

	store := NewElasticsearchStore(client)
	fmt.Println("DEBUG ES: before EnsureIndexes")
	if err := store.EnsureIndexes(ctx); err != nil {
		return err
	}
	fmt.Println("DEBUG ES: after EnsureIndexes")

	ESStore = store
	fmt.Println("DEBUG ES: init complete")
	return nil
}

func (s *ElasticsearchStore) EnsureIndexes(ctx context.Context) error {
	indexMappings := map[string]string{
		IndexSecurityEvents: `{
			"mappings": {
				"properties": {
					"id":          { "type": "keyword" },
					"timestamp":   { "type": "date" },
					"event_type":  { "type": "keyword" },
					"source_ip":   { "type": "keyword" },
					"protocol":    { "type": "keyword" },
					"metadata":    { "type": "object", "enabled": true },
					"ingested_at": { "type": "date" }
				}
			}
		}`,
		IndexAlerts: `{
			"mappings": {
				"properties": {
					"id":           { "type": "keyword" },
					"timestamp":    { "type": "date" },
					"type":         { "type": "keyword" },
					"severity":     { "type": "keyword" },
					"source_ip":    { "type": "keyword" },
					"target":       { "type": "keyword" },
					"description":  { "type": "text" },
					"threat_score": { "type": "float" },
					"status":       { "type": "keyword" },
					"metadata":     { "type": "object", "enabled": true },
					"ingested_at":  { "type": "date" }
				}
			}
		}`,
		IndexResponseActions: `{
			"mappings": {
				"properties": {
					"id":           { "type": "keyword" },
					"alert_id":     { "type": "keyword" },
					"timestamp":    { "type": "date" },
					"action_type":  { "type": "keyword" },
					"source_ip":    { "type": "keyword" },
					"target":       { "type": "keyword" },
					"severity":     { "type": "keyword" },
					"threat_score": { "type": "float" },
					"reason":       { "type": "text" },
					"status":       { "type": "keyword" },
					"metadata":     { "type": "object", "enabled": true },
					"ingested_at":  { "type": "date" }
				}
			}
		}`,
	}

	for index, mapping := range indexMappings {
		if err := s.ensureIndex(ctx, index, mapping); err != nil {
			return err
		}
	}

	return nil
}

func (s *ElasticsearchStore) ensureIndex(ctx context.Context, index string, mapping string) error {
	fmt.Println("DEBUG ES: checking index", index)

	res, err := s.client.Raw().Indices.Exists(
		[]string{index},
		s.client.Raw().Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("check index exists %s: %w", index, err)
	}
	defer res.Body.Close()

	fmt.Println("DEBUG ES: exists status for", index, "=", res.StatusCode)

	if res.StatusCode == 200 {
		fmt.Println("DEBUG ES: index already exists", index)
		return nil
	}

	fmt.Println("DEBUG ES: creating index", index)
	createRes, err := s.client.Raw().Indices.Create(
		index,
		s.client.Raw().Indices.Create.WithContext(ctx),
		s.client.Raw().Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
	)
	if err != nil {
		return fmt.Errorf("create index %s: %w", index, err)
	}
	defer createRes.Body.Close()

	fmt.Println("DEBUG ES: create status for", index, "=", createRes.StatusCode)

	if createRes.IsError() {
		body, _ := io.ReadAll(createRes.Body)
		return fmt.Errorf("create index %s failed: %s", index, string(body))
	}

	fmt.Println("DEBUG ES: created index", index)
	return nil
}

func (s *ElasticsearchStore) IndexDocument(ctx context.Context, index, documentID string, doc interface{}) error {
	fmt.Println("DEBUG ES: indexing document", documentID, "into", index)

	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal document for %s: %w", index, err)
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client.Raw())
	if err != nil {
		return fmt.Errorf("index document into %s: %w", index, err)
	}
	defer res.Body.Close()

	fmt.Println("DEBUG ES: index response status for", index, "=", res.StatusCode)

	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return fmt.Errorf("index document into %s failed: %s", index, string(data))
	}

	return nil
}

func (s *ElasticsearchStore) SearchBySourceIP(ctx context.Context, index, sourceIP string) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"source_ip": sourceIP,
			},
		},
	}

	body, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal search query: %w", err)
	}

	res, err := s.client.Raw().Search(
		s.client.Raw().Search.WithContext(ctx),
		s.client.Raw().Search.WithIndex(index),
		s.client.Raw().Search.WithBody(bytes.NewReader(body)),
		s.client.Raw().Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("search index %s: %w", index, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		data, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("search index %s failed: %s", index, string(data))
	}

	var parsed struct {
		Hits struct {
			Hits []struct {
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("decode search response: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(parsed.Hits.Hits))
	for _, hit := range parsed.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, nil
}

func IndexSecurityEventDoc(doc map[string]interface{}, documentID string) {
    if ESStore == nil {
        return
    }

    doc["ingested_at"] = time.Now().UTC().Format(time.RFC3339Nano)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := ESStore.IndexDocument(ctx, IndexSecurityEvents, documentID, doc); err != nil {
        fmt.Println("failed to index security event in elasticsearch:", err)
    }
}

func IndexAlertDoc(doc map[string]interface{}, documentID string) {
    if ESStore == nil {
        return
    }

    doc["ingested_at"] = time.Now().UTC().Format(time.RFC3339Nano)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := ESStore.IndexDocument(ctx, IndexAlerts, documentID, doc); err != nil {
        fmt.Println("failed to index alert in elasticsearch:", err)
    }
}

func IndexResponseActionDoc(doc map[string]interface{}, documentID string) {
    if ESStore == nil {
        return
    }

    doc["ingested_at"] = time.Now().UTC().Format(time.RFC3339Nano)

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    if err := ESStore.IndexDocument(ctx, IndexResponseActions, documentID, doc); err != nil {
        fmt.Println("failed to index response action in elasticsearch:", err)
    }
}
