package collector

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"sentinelx/models"
	"sentinelx/pipeline"
)

// StartContainerMonitor collects Docker container metrics
func StartContainerMonitor() {

	for {

		cmd := exec.Command(
			"docker",
			"stats",
			"--no-stream",
			"--format",
			"{{.Container}} {{.CPUPerc}} {{.MemUsage}}",
		)

		var out bytes.Buffer
		cmd.Stdout = &out

		err := cmd.Run()
		if err != nil {
			fmt.Println("Container monitor error:", err)
			time.Sleep(10 * time.Second)
			continue
		}

		lines := strings.Split(out.String(), "\n")

		for _, line := range lines {

			fields := strings.Fields(line)

			if len(fields) < 3 {
				continue
			}

			containerID := fields[0]
			cpu := fields[1]
			mem := fields[2]

			// Create SentinelX event
			event := models.NewSecurityEvent("container_metrics")

			event.Metadata["container_id"] = containerID
			event.Metadata["cpu_usage"] = cpu
			event.Metadata["memory_usage"] = mem

			// Debug print
			jsonData, err := event.ToJSON()
			if err == nil {
				fmt.Println(string(jsonData))
			}

			// Send event to pipeline
			pipeline.PublishEvent(event)
		}

		time.Sleep(10 * time.Second)
	}
}
