package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticsearchClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchClient() (*ElasticsearchClient, error) {

	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		return nil, fmt.Errorf("ELASTICSEARCH_URL not set")
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{esURL}, // ✅ from env
	})
	if err != nil {
		return nil, fmt.Errorf("create elasticsearch client: %w", err)
	}

	return &ElasticsearchClient{
		client: es,
	}, nil
}

func (e *ElasticsearchClient) Ping(ctx context.Context) error {
	res, err := e.client.Info(
		e.client.Info.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("elasticsearch ping failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch info returned error: %s", res.String())
	}

	return nil
}

func (e *ElasticsearchClient) Raw() *elasticsearch.Client {
	return e.client
}

func (e *ElasticsearchClient) Close() error {
	return nil
}