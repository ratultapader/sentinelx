package response

import (
	"fmt"
	"net"
	"sync"
)

type RateLimiter struct {
	mu           sync.Mutex
	limitedIPs   map[string]bool
	simulateMode bool
}

func NewRateLimiter(simulate bool) *RateLimiter {
	return &RateLimiter{
		limitedIPs:   make(map[string]bool),
		simulateMode: simulate,
	}
}

func (r *RateLimiter) ValidateIP(ip string) error {
	if ip == "" {
		return fmt.Errorf("missing source ip")
	}
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid ip: %s", ip)
	}
	return nil
}

func (r *RateLimiter) IsProtectedIP(ip string) bool {
	protected := map[string]struct{}{
		"127.0.0.1": {},
		"::1":       {},
	}

	_, exists := protected[ip]
	return exists
}

func (r *RateLimiter) IsLimited(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.limitedIPs[ip]
}

func (r *RateLimiter) ApplyLimit(ip string, limitPerSec int, burst int) (string, error) {
	if err := r.ValidateIP(ip); err != nil {
		return "", err
	}
	if limitPerSec <= 0 {
		return "", fmt.Errorf("invalid limit per second: %d", limitPerSec)
	}
	if burst <= 0 {
		return "", fmt.Errorf("invalid burst: %d", burst)
	}

	if r.IsProtectedIP(ip) {
		return "source ip is protected and cannot be rate limited automatically", nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.limitedIPs[ip] {
		return "rate limit already exists for source ip", nil
	}

	if r.simulateMode {
		r.limitedIPs[ip] = true
		return "simulated rate limit applied", nil
	}

	r.limitedIPs[ip] = true
	return "rate limit applied", nil
}