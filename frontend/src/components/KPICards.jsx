import React from "react";

export default function KPICards({ stats }) {
  return (
    <div style={styles.container}>
      <div style={styles.card}>
        <h4>⏱ MTTR</h4>
        <p>{stats.mttr || 0} min</p>
      </div>

      <div style={styles.card}>
        <h4>🚨 Total Alerts</h4>
        <p>{stats.total_alerts || 0}</p>
      </div>

      <div style={styles.card}>
        <h4>⚠ False Positives</h4>
        <p>{stats.false_positive_rate || 0}%</p>
      </div>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    gap: "15px",
    marginBottom: "20px"
  },
  card: {
    flex: 1,
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    textAlign: "center",
    color: "#e2e8f0",
    boxShadow: "0 6px 20px rgba(0,0,0,0.4)"
  }
};