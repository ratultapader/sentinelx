import React from "react";

export default function AttackPattern({ steps }) {
  if (!steps || steps.length === 0) {
    return (
      <div style={styles.empty}>
        No attack pattern detected
      </div>
    );
  }

  return (
    <div style={styles.container}>
      <h3 style={styles.title}>🚨 Multi-Stage Attack</h3>

      {steps.map((s, i) => (
        <div key={i} style={styles.stepWrapper}>

          {/* STEP NUMBER */}
          <div style={styles.number}>{i + 1}</div>

          {/* STEP CONTENT */}
          <div style={styles.step}>
            <div style={styles.stepTitle}>
              {formatStep(s.step)}
            </div>

            <div style={styles.time}>
              {formatTime(s.timestamp)}
            </div>
          </div>

          {/* ARROW */}
          {i < steps.length - 1 && (
            <div style={styles.arrow}>↓</div>
          )}
        </div>
      ))}
    </div>
  );
}

// ===============================
// HELPERS
// ===============================
function formatStep(step) {
  if (!step) return "Unknown Attack";

  return step
    .replace(/_/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

function formatTime(ts) {
  if (!ts) return "";

  const date = new Date(ts);
  return date.toLocaleTimeString();
}

// ===============================
// STYLES
// ===============================
const styles = {
  container: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
  },
  title: {
    marginBottom: "15px",
    color: "#e2e8f0",
  },
  stepWrapper: {
    display: "flex",
    alignItems: "center",
    marginBottom: "10px",
  },
  number: {
    width: "30px",
    height: "30px",
    borderRadius: "50%",
    background: "#3b82f6",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    marginRight: "10px",
    color: "white",
    fontWeight: "bold",
  },
  step: {
    flex: 1,
    background: "#020617",
    padding: "10px",
    borderRadius: "8px",
  },
  stepTitle: {
    fontWeight: "600",
    color: "#f8fafc",
  },
  time: {
    fontSize: "12px",
    color: "#94a3b8",
    marginTop: "4px",
  },
  arrow: {
    margin: "0 auto",
    color: "#64748b",
    fontSize: "18px",
  },
  empty: {
    color: "#94a3b8",
  },
};