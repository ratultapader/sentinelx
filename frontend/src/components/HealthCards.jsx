import React from "react";

export default function HealthCards({ health }) {
  return (
    <div style={styles.container}>
      <Card name="Database" status={health.database} />
      <Card name="Elasticsearch" status={health.elasticsearch} />
      <Card name="Neo4j" status={health.neo4j} />
    </div>
  );
}

function Card({ name, status }) {
  const ok = status === "ok";

  return (
    <div
      style={{
        ...styles.card,
        background: ok ? "#064e3b" : "#7f1d1d"
      }}
    >
      <h4>{name}</h4>
      <p>{ok ? "✅ Healthy" : "❌ Down"}</p>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    gap: "20px"
  },
  card: {
    padding: "15px",
    borderRadius: "10px",
    width: "200px",
    color: "white"
  }
};