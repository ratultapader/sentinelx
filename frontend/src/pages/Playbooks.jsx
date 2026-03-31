import React, { useEffect, useState } from "react";
import PlaybookTable from "../components/PlaybookTable";

export default function Playbooks() {
  const [playbooks, setPlaybooks] = useState([]);
  const [score, setScore] = useState(0.8);
  const [action, setAction] = useState("block_ip");

  const load = () => {
    fetch("http://localhost:9090/api/playbooks", {
      headers: { "X-Tenant-ID": "t1" }
    })
      .then(res => res.json())
      .then(setPlaybooks);
  };

  useEffect(() => {
    load();
  }, []);

  const addPlaybook = () => {
    const condition = `threat_score > ${score}`;

    // ✅ Prevent duplicate
    if (playbooks.some(p => p.condition === condition && p.action === action)) {
      alert("⚠️ Duplicate playbook not allowed");
      return;
    }

    fetch("http://localhost:9090/api/playbooks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "t1"
      },
      body: JSON.stringify({ condition, action })
    }).then(load);
  };

  return (
    <div style={styles.container}>
      <h2>⚙️ Response Playbooks</h2>

      {/* 🔥 INPUT PANEL */}
      <div style={styles.controls}>
        <label>Threat Score:</label>
        <input
          type="number"
          step="0.1"
          value={score}
          onChange={(e) => setScore(e.target.value)}
          style={styles.input}
        />

        <label>Action:</label>
        <select
          value={action}
          onChange={(e) => setAction(e.target.value)}
          style={styles.input}
        >
          <option value="block_ip">block_ip</option>
          <option value="rate_limit">rate_limit</option>
          <option value="alert">alert</option>
        </select>

        <button onClick={addPlaybook} style={styles.button}>
          + Add Playbook
        </button>
      </div>

      <PlaybookTable playbooks={playbooks} onRefresh={load} />
    </div>
  );
}

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },
  controls: {
    display: "flex",
    gap: "10px",
    marginBottom: "15px",
    alignItems: "center"
  },
  input: {
    padding: "6px",
    background: "#020617",
    color: "white",
    border: "1px solid #334155"
  },
  button: {
    padding: "6px 12px",
    background: "#2563eb",
    color: "white",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer"
  }
};