import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import TimelineView from "../components/TimelineView";
import GraphView from "../components/GraphView";
import AttackPattern from "../components/AttackPattern";

export default function Investigation() {
  const query = new URLSearchParams(useLocation().search);
  const navigate = useNavigate();

  const ip = query.get("ip");

  const [attackSteps, setAttackSteps] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // ===============================
  // 🔒 BLOCK DIRECT ACCESS
  // ===============================
  useEffect(() => {
    if (!ip) navigate("/incidents");
  }, [ip, navigate]);

  // ===============================
  // 🔥 FETCH ATTACK PATTERN
  // ===============================
  useEffect(() => {
    if (!ip) return;

    setLoading(true);
    setError(null);

    fetch(`http://localhost:9090/api/attack_pattern/${ip}`, {
      headers: { "X-Tenant-ID": "t1" }
    })
      .then(res => {
        if (!res.ok) throw new Error("Failed to load attack pattern");
        return res.json();
      })
      .then(data => setAttackSteps(Array.isArray(data) ? data : []))
      .catch(err => {
        console.error("AttackPattern error:", err);
        setError("Failed to load attack pattern");
        setAttackSteps([]);
      })
      .finally(() => setLoading(false));
  }, [ip]);

  if (!ip) return null;

  return (
    <div style={styles.container}>

      {/* ================= HEADER ================= */}
      <div style={styles.header}>
        <button onClick={() => navigate("/incidents")} style={styles.backBtn}>
          ⬅ Back
        </button>

        <div>
          <h2 style={{ margin: 0 }}>🔍 Investigation</h2>
          <span style={styles.ip}>{ip}</span>
        </div>

        <div style={styles.status}>
          {loading ? "⏳ Loading..." : "🟢 Active"}
        </div>
      </div>

      {/* ================= ATTACK PATTERN ================= */}
      <div style={styles.section}>
        <h3 style={styles.title}>⚔️ Attack Pattern</h3>

        {loading && <p style={styles.muted}>Loading attack pattern...</p>}
        {error && <p style={styles.error}>{error}</p>}

        {!loading && !error && attackSteps.length === 0 && (
          <p style={styles.muted}>No attack pattern found</p>
        )}

        {!loading && !error && attackSteps.length > 0 && (
          <AttackPattern steps={attackSteps} />
        )}
      </div>

      {/* ================= MAIN ================= */}
      <div style={styles.main}>

        {/* LEFT: TIMELINE */}
        <div style={styles.left}>
          <h3 style={styles.title}>📜 Timeline</h3>

          {/* ✅ FIX: Timeline now handled inside component */}
          <TimelineView ip={ip} />
        </div>

        {/* RIGHT: GRAPH */}
        <div style={styles.right}>
          <h3 style={styles.title}>🌐 Attack Graph</h3>
          <GraphView ip={ip} />
        </div>

      </div>
    </div>
  );
}

// ===============================
// 🎨 STYLES
// ===============================
const styles = {
  container: {
    background: "#020617",
    minHeight: "100vh",
    color: "#e2e8f0"
  },
  header: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    padding: "15px",
    borderBottom: "1px solid #1e293b"
  },
  backBtn: {
    background: "#1e293b",
    border: "none",
    color: "white",
    padding: "6px 12px",
    borderRadius: "6px",
    cursor: "pointer"
  },
  ip: {
    fontSize: "12px",
    color: "#94a3b8"
  },
  status: {
    background: "#022c22",
    color: "#4ade80",
    padding: "5px 10px",
    borderRadius: "6px",
    fontSize: "12px"
  },
  section: {
    padding: "15px",
    borderBottom: "1px solid #1e293b"
  },
  main: {
    display: "flex",
    height: "calc(100vh - 200px)"
  },
  left: {
    width: "40%",
    padding: "15px",
    overflowY: "auto",
    borderRight: "1px solid #1e293b"
  },
  right: {
    width: "60%",
    padding: "15px",
    overflow: "hidden"
  },
  title: {
    marginBottom: "10px"
  },
  muted: {
    color: "#94a3b8"
  },
  error: {
    color: "#f87171"
  }
};