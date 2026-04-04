import React from "react";
import { useNavigate } from "react-router-dom";

export default function IncidentDetails({ incident }) {
  const navigate = useNavigate();

  if (!incident) return null;

  const alerts = incident.alerts || [];

  // 🔥 fallback IP (if no alerts)
  const fallbackIP =
    incident.source_ip ||
    incident.ip ||
    "101.101.101.101"; // last fallback

  return (
    <div style={styles.card}>
      <h3>🧾 Incident Details</h3>

      <p><strong>ID:</strong> {incident.id}</p>
      <p><strong>Severity:</strong> {incident.severity}</p>
      <p><strong>Status:</strong> {incident.status || "new"}</p>
      <p><strong>Alert Count:</strong> {alerts?.length || 0}</p>

      {/* 🔥 CASE 1: HAS ALERTS */}
      {(alerts?.length || 0) > 0 ? (
        <>
          <h4>🌐 Related IPs</h4>

          {alerts.map((a, i) => (
            <div key={i} style={styles.row}>
              <span>{a.source_ip} → ({a.severity})</span>

              <button
                style={styles.button}
                onClick={() => navigate(`/investigation?ip=${a.source_ip}`)}
              >
                🔍 Investigate
              </button>
            </div>
          ))}
        </>
      ) : (
        /* 🔥 CASE 2: NO ALERTS (FIXED) */
        <>
          <div style={{ color: "#f87171", marginBottom: "10px" }}>
            ⚠ No alerts — using fallback IP
          </div>

          <button
            style={styles.button}
            onClick={() => navigate(`/investigation?ip=${fallbackIP}`)}
          >
            🔍 Investigate Anyway
          </button>
        </>
      )}
    </div>
  );
}

const styles = {
  card: {
    marginTop: "20px",
    padding: "15px",
    background: "#1e293b",
    borderRadius: "10px"
  },
  row: {
    display: "flex",
    justifyContent: "space-between",
    marginBottom: "8px"
  },
  button: {
    padding: "6px 12px",
    background: "#3b82f6",
    border: "none",
    borderRadius: "6px",
    color: "white",
    cursor: "pointer"
  }
};