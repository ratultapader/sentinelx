import React from "react";

export default function AlertTable({ alerts, onSelect }) {
  return (
    <table width="100%" style={{ background: "#1e293b", color: "white" }}>
      <thead>
        <tr>
          <th>Severity</th>
          <th>Source IP</th>
          <th>Target</th>
          <th>Score</th>
          <th>Time</th>
        </tr>
      </thead>

      <tbody>
        {alerts.map((a, i) => (
          <tr key={i} onClick={() => onSelect(a)} style={{ cursor: "pointer" }}>
            <td style={{ color: getColor(a.severity) }}>{a.severity}</td>
            <td>{a.source_ip}</td>
            <td>{a.destination}</td>
            <td>{a.threat_score}</td>
            <td>{new Date(a.timestamp).toLocaleString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function getColor(severity) {
  switch (severity) {
    case "critical": return "#ef4444";
    case "high": return "#f97316";
    case "medium": return "#eab308";
    default: return "#22c55e";
  }
}