import React, { useEffect, useState } from "react";
import StatsCards from "../components/ThreatStats";
import ThreatTable from "../components/ThreatTable";
import api from "../services/api"; // ✅ use axios

export default function ThreatIntel() {
  const [ips, setIps] = useState([]);
  const [newIP, setNewIP] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const loadData = async () => {
    setLoading(true);
    setError(null);

    try {
      const res = await api.get("/api/threat_intel"); // ✅ tenant auto

      const sorted = (res.data || []).sort(
        (a, b) => (b.match_count || 0) - (a.match_count || 0)
      );

      setIps(sorted);
      setLoading(false);

    } catch (err) {
      console.error(err);
      setError("Failed to load threat intel");
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();

    const interval = setInterval(loadData, 5000); // 🔄 auto refresh
    return () => clearInterval(interval);
  }, []);

  const addIP = async () => {
    if (!newIP) return;

    try {
      await api.post("/api/threat_intel", { ip: newIP }); // ✅ clean

      setNewIP("");
      loadData();

    } catch (err) {
      console.error(err);
    }
  };

  return (
    <div style={styles.container}>
      <h2 style={styles.title}>🧠 Threat Intelligence</h2>

      {/* ADD IP */}
      <div style={styles.inputBox}>
        <input
          placeholder="Add malicious IP..."
          value={newIP}
          onChange={(e) => setNewIP(e.target.value)}
          style={styles.input}
        />
        <button onClick={addIP} style={styles.button}>
          Add
        </button>
      </div>

      {/* STATES */}
      {loading && <p>Loading threat intelligence...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      {/* STATS */}
      <StatsCards data={ips} />

      {/* TABLE */}
      <ThreatTable ips={ips} />
    </div>
  );
}

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },
  title: {
    marginBottom: "20px"
  },
  inputBox: {
    display: "flex",
    gap: "10px",
    marginBottom: "20px"
  },
  input: {
    padding: "10px",
    borderRadius: "6px",
    border: "none",
    width: "250px"
  },
  button: {
    padding: "10px 15px",
    border: "none",
    borderRadius: "6px",
    background: "#3b82f6",
    color: "white",
    cursor: "pointer"
  }
};