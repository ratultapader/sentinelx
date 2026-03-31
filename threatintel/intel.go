package threatintel

// 🔥 store match counts
var matchCounts = make(map[string]int)

// ===============================
// EXISTING CHECK
// ===============================
func IsMaliciousIP(ip string) bool {
	_, exists := maliciousIPs[ip]
	return exists
}

// ===============================
// NEW: INCREMENT MATCH
// ===============================
func IncrementMatch(ip string) {
	matchCounts[ip]++
}

// ===============================
// GET ALL THREATS (FIXED)
// ===============================
func GetAllThreats() []map[string]interface{} {
	var result []map[string]interface{}

	for ip := range maliciousIPs {
		result = append(result, map[string]interface{}{
			"ip":          ip,
			"reason":      "Threat Feed",
			"source":      "File",
			"match_count": matchCounts[ip], // 🔥 FIXED
		})
	}

	return result
}

// ===============================
// ADD THREAT
// ===============================
func AddThreat(ip string) {
	maliciousIPs[ip] = true
}