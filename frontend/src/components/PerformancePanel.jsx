import React, { useEffect, useState } from "react";
import api from "../services/api";

export default function PerformancePanel() {
  const [metrics, setMetrics] = useState(null);

  useEffect(() => {
    api.get("/api/performance").then(res => setMetrics(res.data));
  }, []);

  if (!metrics) return <p>Loading performance...</p>;

  return (
    <div style={styles.card}>
      <h3>⚡ Performance Metrics</h3>

      <div style={styles.grid}>
        <div style={styles.metric}>
          <h4>Events/sec</h4>
          <p>{metrics.events_per_sec}</p>
        </div>

        <div style={styles.metric}>
          <h4>Alerts/sec</h4>
          <p>{metrics.alerts_per_sec}</p>
        </div>

        <div style={styles.metric}>
          <h4>Latency</h4>
          <p>{metrics.latency} ms</p>
        </div>
      </div>
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    marginBottom: "20px",
    borderRadius: "10px",
    color: "white"
  },
  grid: {
    display: "flex",
    justifyContent: "space-around",
    marginTop: "10px"
  },
  metric: {
    textAlign: "center"
  }
};