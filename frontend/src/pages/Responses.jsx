import React, { useEffect, useState } from "react";
import api from "../services/api";

export default function Responses() {
  const [actions, setActions] = useState([]);

  useEffect(() => {
    api
      .get("/api/response_actions", {
        headers: { "X-Tenant-ID": "t1" },
      })
      .then((res) => setActions(res.data || []))
      .catch(() => setActions([]));
  }, []);

  return (
    <div className="container">
      <h1>⚡ Response Actions</h1>

      {actions.length === 0 ? (
        <p>No actions found</p>
      ) : (
        actions.map((a, i) => (
          <div key={i} style={styles.card}>
            <strong>{a.type}</strong> — {a.status}
            <div>Target: {a.target}</div>
          </div>
        ))
      )}
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "10px",
    marginBottom: "10px",
    borderRadius: "8px",
    color: "#e2e8f0",
  },
};