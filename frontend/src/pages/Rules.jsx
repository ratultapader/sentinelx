import React, { useEffect, useState } from "react";
import RuleForm from "../components/RuleForm";
import RuleTable from "../components/RuleTable";
import { get } from "../services/api";

export default function Rules() {
  const [rules, setRules] = useState([]);
  const [showForm, setShowForm] = useState(false);

  const loadRules = async () => {
    try {
      const data = await get("/api/rules");
      setRules(data);
    } catch (err) {
      console.error("Failed to load rules", err);
    }
  };

  useEffect(() => {
    loadRules();
  }, []);

  return (
    <div style={styles.container}>
      
      {/* 🔥 HEADER */}
      <div style={styles.header}>
        <h1 style={styles.title}>⚙️ Rule Engine</h1>
        <button style={styles.button} onClick={() => setShowForm(!showForm)}>
          + Add Rule
        </button>
      </div>

      {/* 🔥 TOTAL COUNT */}
      <h3 style={styles.count}>Total Rules: {rules.length}</h3>

      {/* 🔥 FORM */}
      {showForm && <RuleForm onCreated={loadRules} />}

      {/* 🔥 TABLE */}
      <RuleTable rules={rules} onToggle={loadRules} />

    </div>
  );
}

/* ================= STYLES ================= */

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },
  header: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: "10px"
  },
  title: {
    margin: 0
  },
  button: {
    background: "#3b82f6",
    color: "white",
    border: "none",
    padding: "8px 14px",
    borderRadius: "6px",
    cursor: "pointer",
    fontWeight: "500"
  },
  count: {
    marginTop: "10px",
    marginBottom: "15px",
    color: "#94a3b8"
  }
};