import React from "react";

function RecentAlerts({ alerts = [], summary }) {

  const getSeverityColor = (severity) => {
    switch (severity) {
      case "critical": return "#ef4444"; // red
      case "high": return "#f97316";     // orange
      case "medium": return "#eab308";   // yellow
      case "low": return "#22c55e";      // green
      default: return "#94a3b8";         // gray
    }
  };

  const formatTime = (ts) => {
    if (!ts) return "";
    return new Date(ts).toLocaleString();
  };

  return (
    <div style={{ marginTop: "20px" }}>
      <h3>Recent Alerts</h3>

      {alerts.length > 0 ? (
        alerts.map((alert, index) => (
          <div
            key={index}
            style={{
              background: "#1e293b",
              padding: "12px",
              marginBottom: "12px",
              borderRadius: "10px",
              borderLeft: `5px solid ${getSeverityColor(alert.severity)}`
            }}
          >
            <div><b>IP:</b> {alert.source_ip}</div>

            <div>
              <b>Severity:</b>{" "}
              <span style={{ color: getSeverityColor(alert.severity), fontWeight: "bold" }}>
                {alert.severity}
              </span>
            </div>

            <div><b>Score:</b> {alert.threat_score}</div>

            <div><b>Target:</b> {alert.destination || alert.target}</div>

            <div style={{ fontSize: "12px", color: "#94a3b8" }}>
              {formatTime(alert.timestamp)}
            </div>
          </div>
        ))
      ) : (
        <div style={{ opacity: 0.8 }}>
          <p>No detailed alerts yet</p>
          <p>Total Alerts: {summary?.total_alerts || 0}</p>
        </div>
      )}
    </div>
  );
}

export default RecentAlerts;