import React from "react";

export default function ThreatScoreCard({ score = 0, breakdown = {} }) {

  // 🔥 Color based on severity
  const color =
    score > 0.8 ? "#ef4444" :   // red (critical)
    score > 0.6 ? "#f59e0b" :   // orange (medium)
    "#22c55e";                  // green (low)

  // 🔥 SOC-style labels
  const labels = {
    anomaly: "Anomaly Detection",
    signature: "Signature Match",
    reputation: "IP Reputation",
    behavior: "Behavior Analysis"
  };

  return (
    <div style={styles.card}>
      
      {/* 🔥 Header */}
      <h3 style={{ color }}>
        🔥 Threat Score: {score.toFixed(2)}
      </h3>

      {/* 🔥 Breakdown Bars */}
      {Object.entries(breakdown).map(([key, value]) => (
        <Bar
          key={key}
          label={labels[key] || key}
          value={value ?? 0}
        />
      ))}

    </div>
  );
}

/* ================= BAR COMPONENT ================= */

function Bar({ label, value }) {
  return (
    <div style={styles.barContainer}>
      
      {/* Label */}
      <div style={styles.label}>
        {label} ({value.toFixed(2)})
      </div>

      {/* Progress bar */}
      <div style={styles.barBg}>
        <div
          style={{
            ...styles.barFill,
            width: `${Math.min(value * 100, 100)}%`
          }}
        />
      </div>

    </div>
  );
}

/* ================= STYLES ================= */

const styles = {
  card: {
    marginTop: "15px",
    background: "#020617",
    padding: "15px",
    borderRadius: "10px",
    color: "#e2e8f0",
    boxShadow: "0 4px 12px rgba(0,0,0,0.4)"
  },

  barContainer: {
    marginBottom: "12px"
  },

  label: {
    fontSize: "13px",
    marginBottom: "4px",
    color: "#cbd5f5"
  },

  barBg: {
    height: "8px",
    background: "#334155",
    borderRadius: "5px",
    overflow: "hidden"
  },

  barFill: {
    height: "8px",
    background: "#3b82f6",
    borderRadius: "5px",
    transition: "width 0.5s ease"
  }
};