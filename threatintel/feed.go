package threatintel

import (
	"bufio"
	"fmt"
	"os"
)

var maliciousIPs = make(map[string]bool)

func LoadThreatFeed(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		ip := scanner.Text()

		if ip != "" {
			maliciousIPs[ip] = true
		}
	}

	fmt.Printf("Threat intel loaded: %d malicious IPs\n", len(maliciousIPs))

	return nil
}