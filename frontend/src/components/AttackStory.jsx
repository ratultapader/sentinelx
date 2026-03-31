import React from "react";

export default function AttackStory({ story }) {
  return (
    <div style={styles.card}>
      <h3>🧠 Attack Story</h3>
      <p>{story}</p>
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    color: "#e2e8f0",
    boxShadow: "0 6px 20px rgba(0,0,0,0.4)"
  }
};