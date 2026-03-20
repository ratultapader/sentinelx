package threatfeed

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	threatIPs = make(map[string]bool)
	mutex     sync.RWMutex
)

func IsMalicious(ip string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return threatIPs[ip]
}

func AddTestIP(ip string) {
	mutex.Lock()
	defer mutex.Unlock()
	threatIPs[ip] = true
}

func Count() int {
	mutex.RLock()
	defer mutex.RUnlock()
	return len(threatIPs)
}

func loadFeed(client *http.Client, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	newIPs := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		ip := fields[0]
		if net.ParseIP(ip) != nil {
			newIPs[ip] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	mutex.Lock()
	for ip := range newIPs {
		threatIPs[ip] = true
	}
	mutex.Unlock()

	return nil
}

func StartThreatFeedUpdater() {
	feeds := []string{
		"https://feodotracker.abuse.ch/downloads/ipblocklist.txt",
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	for _, feed := range feeds {
		_ = loadFeed(client, feed)
	}

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			for _, feed := range feeds {
				_ = loadFeed(client, feed)
			}
		}
	}()
}