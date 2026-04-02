import React from "react";

export default function FeedbackPanel({ onFeedback }) {
  return (
    <div style={styles.card}>
      <h3>🧠 Analyst Feedback</h3>

      <button onClick={() => onFeedback("true_positive")}>
        ✅ True Positive
      </button>

      <button onClick={() => onFeedback("false_positive")}>
        ❌ False Positive
      </button>
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    color: "#e2e8f0"
  }
};