import React, { useEffect, useState } from "react";
import api from "../services/api";

export default function AuditLogs() {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    api.get("/api/audit_logs").then(res => setLogs(res.data));
  }, []);

  return (
    <div style={styles.card}>
      <h3>📜 Audit Logs</h3>

      {logs.length === 0 && <p>No logs available</p>}

      {logs.map((l, i) => {
        // 🔥 detect high alert count
        const isCritical = l.action.includes("40") || l.action.includes("30");

        return (
          <div key={i} style={styles.row}>
            
            {/* ⏱ TIME */}
            <span style={styles.time}>
              {new Date(l.timestamp).toLocaleTimeString()}
            </span>

            {/* 🔥 ACTION */}
            <span
              style={{
                ...styles.action,
                color: isCritical ? "#f87171" : "#e2e8f0"
              }}
            >
              {l.user} → {l.action} → <b>{l.target}</b>
            </span>

          </div>
        );
      })}
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    marginBottom: "20px",
    borderRadius: "10px",
    color: "white"
  },
  row: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: "8px",
    borderBottom: "1px solid #334155",
    paddingBottom: "5px"
  },
  time: {
    color: "#94a3b8",
    fontSize: "12px",
    minWidth: "80px"
  },
  action: {
    fontSize: "14px"
  }
};