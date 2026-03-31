package configs

import "os"

type Config struct {
	Port             string
	ElasticsearchURL string
	Neo4jURI         string
	Neo4jUsername    string
	Neo4jPassword    string
}

func Load() Config {
	return Config{
		Port:             getEnv("PORT", "9090"),
		ElasticsearchURL: getEnv("ES_URL", "http://localhost:9200"),
		Neo4jURI:         getEnv("NEO4J_URI", "neo4j://localhost:7687"),
		Neo4jUsername:    getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword:    getEnv("NEO4J_PASSWORD", "password"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}