package app

import (
	"fmt"

	"sentinelx/configs"
)

// Bootstrap initializes full application
func Bootstrap() *Dependencies {

	cfg := configs.Load()

	fmt.Println("🚀 Bootstrapping SentinelX")
	fmt.Println("Port:", cfg.Port)

	deps := InitDependencies()

	return deps
}