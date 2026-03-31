import React from "react";

export default function IncidentTable({ incidents, onSelect }) {
  // ✅ SAFETY: ensure it's always an array
  const safeIncidents = Array.isArray(incidents) ? incidents : [];

  return (
    <table width="100%">
      <thead>
        <tr>
          <th>ID</th>
          <th>Severity</th>
          <th>Status</th>
          <th>Alerts</th>
          <th>Time</th>
        </tr>
      </thead>

      <tbody>
        {safeIncidents.length === 0 ? (
          <tr>
            <td colSpan="5" style={{ textAlign: "center", padding: "10px" }}>
              No incidents found
            </td>
          </tr>
        ) : (
          safeIncidents.map((i, index) => (
            <tr
              key={index}
              onClick={() =>
  onSelect({
    ...i,
    id: i.id || i._id || i.incident_id, // ✅ ensure ID exists
  })
}
              style={{ cursor: "pointer" }}
            >
             <td>{i.id || i._id || i.incident_id}</td>
              <td style={{ color: getColor(i.severity) }}>
                {i.severity}
              </td>
              <td>{i.status}</td>
              <td>{i.alert_count}</td>
              <td>{formatTime(i.timestamp)}</td>
            </tr>
          ))
        )}
      </tbody>
    </table>
  );
}

function getColor(severity) {
  switch (severity) {
    case "critical":
      return "#ef4444";
    case "high":
      return "#f97316";
    case "medium":
      return "#eab308";
    default:
      return "#22c55e";
  }
}

function formatTime(ts) {
  return ts ? new Date(ts).toLocaleString() : "";
}