import React from "react";

export default function AttackStory({ steps = [] }) {

  const icons = ["🚀", "💉", "🧨", "🔓", "⚠️"];

  return (
    <div style={styles.card}>
      <h3>🧠 Attack Story</h3>

      {steps.length === 0 ? (
        <p>No attack steps available</p>
      ) : (
        steps.map((s, i) => (
          <div key={i} style={styles.step}>
            <strong>
              {icons[i % icons.length]} Step {i + 1}:
            </strong>{" "}
            {s}
          </div>
        ))
      )}
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
  },
  step: {
    marginBottom: "10px",
    padding: "6px",
    background: "#0f172a",
    borderRadius: "6px"
  }
};