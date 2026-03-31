import React, { useState } from "react";
import { postNoBody } from "../services/api";

export default function RuleTable({ rules, onToggle }) {

  const [loadingId, setLoadingId] = useState(null);

  const toggleRule = async (rule) => {
    try {
      setLoadingId(rule.id);

      await postNoBody(`/api/rules/${rule.id}/toggle`);

      onToggle();
    } catch (err) {
      console.error("Toggle failed", err);
    } finally {
      setLoadingId(null);
    }
  };

  return (
    <table style={styles.table}>
      <thead>
        <tr>
          <th>Name</th>
          <th>Condition</th>
          <th>Action</th>
          <th>Enabled</th>
        </tr>
      </thead>

      <tbody>
        {rules.map((r) => (
          <tr key={r.id} style={styles.row}>
            <td>{r.name}</td>
            <td>{r.condition}</td>
            <td>{r.action}</td>
            <td>
              <button
                onClick={() => toggleRule(r)}
                disabled={loadingId === r.id}
                style={{
                  ...styles.button,
                  background: r.enabled ? "#22c55e" : "#ef4444", // 🔥 green / red
                  opacity: loadingId === r.id ? 0.6 : 1,
                  cursor: loadingId === r.id ? "not-allowed" : "pointer"
                }}
              >
                {loadingId === r.id
                  ? "..."
                  : r.enabled
                  ? "ON"
                  : "OFF"}
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

/* ================= STYLES ================= */

const styles = {
  table: {
    width: "100%",
    marginTop: "20px",
    borderCollapse: "collapse",
    background: "#1e293b",
    color: "white"
  },
  row: {
    borderBottom: "1px solid #334155"
  },
  button: {
    border: "none",
    color: "white",
    padding: "6px 12px",
    borderRadius: "6px",
    fontWeight: "600",
    transition: "0.2s"
  }
};