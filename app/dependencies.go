package app

import "sentinelx/storage"

// Dependencies holds core services
type Dependencies struct {
	ES    interface{}
	Graph interface{}
	DB    interface{}
}

// InitDependencies initializes core services
func InitDependencies() *Dependencies {
	return &Dependencies{
		ES:    storage.ESStore,
		Graph: storage.GraphStore,
		DB:    storage.DB,
	}
}