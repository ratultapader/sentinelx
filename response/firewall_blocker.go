package response

import (
	"fmt"
	"net"
	"sync"
)

type FirewallBlocker struct {
	mu           sync.Mutex
	blockedIPs   map[string]bool
	simulateMode bool
}

func NewFirewallBlocker(simulate bool) *FirewallBlocker {
	return &FirewallBlocker{
		blockedIPs:   make(map[string]bool),
		simulateMode: simulate,
	}
}

func (b *FirewallBlocker) IsProtectedIP(ip string) bool {
	protected := map[string]struct{}{
		"127.0.0.1": {},
		"::1":       {},
	}

	_, exists := protected[ip]
	return exists
}

func (b *FirewallBlocker) ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("missing source ip")
	}
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid ip: %s", ip)
	}
	return nil
}

func (b *FirewallBlocker) IsBlocked(ip string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.blockedIPs[ip]
}

func (b *FirewallBlocker) BlockIP(ip string) (string, error) {
	if err := b.ValidateIP(ip); err != nil {
		return "", err
	}

	if b.IsProtectedIP(ip) {
		return "source ip is protected and cannot be blocked automatically", nil
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.blockedIPs[ip] {
		return "ip already blocked", nil
	}

	if b.simulateMode {
		b.blockedIPs[ip] = true
		return "simulated firewall block applied", nil
	}

	// Real firewall execution can be added later here.
	// Example future implementation:
	// iptables -A INPUT -s <ip> -j DROP

	b.blockedIPs[ip] = true
	return "firewall block applied", nil
}