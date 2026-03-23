package investigation

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GenerateTextReport(t AttackTimeline) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Attack Investigation Report\n")
	fmt.Fprintf(&b, "Source IP: %s\n", t.SourceIP)
	fmt.Fprintf(&b, "Attack Type: %s\n", t.AttackType)
	fmt.Fprintf(&b, "Risk Level: %s\n", t.RiskLevel)
	fmt.Fprintf(&b, "Start Time: %s\n", t.StartTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "End Time: %s\n", t.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "Event Count: %d\n", t.EventCount)
	fmt.Fprintf(&b, "Stages: %s\n", strings.Join(t.Stages, " -> "))
	fmt.Fprintf(&b, "Conclusion: %s\n", t.Conclusion)
	fmt.Fprintf(&b, "\nTimeline:\n")

	for _, ev := range t.Events {
		fmt.Fprintf(
			&b,
			"- %s | %s | %s | %s | severity=%s score=%.2f\n",
			ev.Timestamp.Format("15:04:05"),
			ev.Stage,
			ev.EventType,
			ev.Summary,
			ev.Severity,
			ev.ThreatScore,
		)
	}

	fmt.Fprintf(&b, "\nRecommendations:\n")
	for _, rec := range t.Recommendations {
		fmt.Fprintf(&b, "- %s\n", rec)
	}

	return b.String()
}

func GenerateJSONReport(t AttackTimeline) ([]byte, error) {
	return json.MarshalIndent(t, "", "  ")
}