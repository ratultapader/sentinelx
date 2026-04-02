import React from "react";

export default function PriorityCard({ score, reasons = [] }) {

  const bg =
    score > 80 ? "#7f1d1d" :
    score > 60 ? "#78350f" :
    "#1e293b";

  return (
    <div style={{ ...styles.card, background: bg }}>
      
      <h3>🔥 Priority Score: {score}/100</h3>

      <h4 style={{ marginTop: "10px" }}>
        Why this alert is critical:
      </h4>

      <ul style={{ marginTop: "8px", paddingLeft: "18px" }}>
        {reasons.map((r, i) => (
          <li key={i} style={{ marginBottom: "6px" }}>
            ✔ {r}
          </li>
        ))}
      </ul>

    </div>
  );
}

const styles = {
  card: {
    padding: "15px",
    borderRadius: "10px",
    color: "white",
    marginBottom: "20px",
    boxShadow: "0 6px 20px rgba(0,0,0,0.4)"
  }
};