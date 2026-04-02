package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 🔹 Config struct (central config)
type Config struct {
	Port             string
	ElasticsearchURL string

	Neo4jURI      string
	Neo4jUsername string
	Neo4jPassword string
}

// 🔹 Load config from .env + system env
func Load() Config {

	// 🔥 Load .env file (optional in prod)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system env")
	}

	cfg := Config{
		Port:             getEnv("PORT", "9090"),
		ElasticsearchURL: getEnv("ELASTICSEARCH_URL", "http://localhost:9200"),

		Neo4jURI:      getEnv("NEO4J_URI", "neo4j://localhost:7687"),
		Neo4jUsername: getEnv("NEO4J_USERNAME", "neo4j"),
		Neo4jPassword: getEnv("NEO4J_PASSWORD", "password"),
	}

	validate(cfg)

	return cfg
}

// 🔹 Helper to read env
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// 🔹 Validate required configs
func validate(cfg Config) {
	if cfg.ElasticsearchURL == "" {
		log.Fatal("ELASTICSEARCH_URL is required")
	}

	if cfg.Neo4jURI == "" {
		log.Fatal("NEO4J_URI is required")
	}
}