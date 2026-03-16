package threatintel

func IsMaliciousIP(ip string) bool {

	_, exists := maliciousIPs[ip]

	return exists
}