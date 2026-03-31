import React from "react";

export default function MetricsPanel({ metrics }) {
  return (
    <div style={styles.container}>
      <Metric title="Total Events" value={metrics.events_processed_total} />
      <Metric title="Total Alerts" value={metrics.alerts_generated_total} />
      <Metric title="Unique Attackers" value={metrics.unique_attackers} />
      <Metric title="Attack Types" value={metrics.attack_types} />
    </div>
  );
}

function Metric({ title, value }) {
  return (
    <div style={styles.card}>
      <h4>{title}</h4>
      <p>{value || 0}</p>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    gap: "20px"
  },
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    width: "200px",
    color: "white"
  }
};