import React from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid
} from "recharts";

export default function ThreatTrendChart({ data }) {
  return (
    <div style={styles.card}>
      <h3>📈 Threat Trends</h3>

      <LineChart width={600} height={300} data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis />
        <Tooltip />

        <Line type="monotone" dataKey="alerts" stroke="#ef4444" />
      </LineChart>
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    color: "#e2e8f0",
    marginBottom: "20px"
  }
};