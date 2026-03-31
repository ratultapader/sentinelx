import React from "react";

export default function StatsCards({ data }) {
  const total = data.length;
  const top = data[0]?.ip || "N/A";
  const matches = data.reduce(
    (sum, i) => sum + (i.match_count || 0),
    0
  );

  return (
    <div style={styles.container}>
      <Card title="Total Threat IPs" value={total} />
      <Card title="Matches Today" value={matches} />
      <Card title="Top Malicious IP" value={top} />
    </div>
  );
}

function Card({ title, value }) {
  return (
    <div style={styles.card}>
      <h4>{title}</h4>
      <p style={styles.value}>{value}</p>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    gap: "20px",
    marginBottom: "20px"
  },
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    width: "220px"
  },
  value: {
    fontSize: "20px",
    fontWeight: "bold"
  }
};