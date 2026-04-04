import React, { useEffect, useState, useRef } from "react";
import { fetchRecentAlerts } from "../services/dashboardApi";
import Filters from "../components/Filters";
import AlertTable from "../components/AlertTable";
import AlertDetails from "../components/AlertDetails";

// 🔥 EXPORT
import { exportToCSV, exportToJSON } from "../utils/export";

export default function Alerts() {
  const [alerts, setAlerts] = useState([]);
  const [filtered, setFiltered] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);

  const userSelectedRef = useRef(false);

  const loadAlerts = async () => {
    try {
      setLoading(true);

      const data = await fetchRecentAlerts();
      console.log("RAW ALERTS:", data);

      let items = Array.isArray(data) ? data : data?.items || [];

      // 🔥 NORMALIZE + INCLUDE SCORING
      const normalized = items.map((item, index) => ({
        id: item.id || item._id || `alert-${index}`,
        type: item.event_type || item.type || "Unknown",
        source_ip: item.source_ip || item.ip || "N/A",
        severity: item.severity || "low",
        timestamp: item.timestamp || new Date().toISOString(),
        target: item.target || "unknown",

        // 🔥 IMPORTANT (SCORING)
        threat_score: item.threat_score || 0,
        anomaly_score: item.anomaly_score || 0,
        signature_score: item.signature_score || 0,
        ip_reputation: item.ip_reputation || 0,
        behavior_score: item.behavior_score || 0,

        ...item
      }));

      setAlerts(normalized);
      setFiltered(normalized);

      if ((normalized?.length || 0) > 0 && !userSelectedRef.current) {
        setSelected(normalized[0]);
      }

    } catch (err) {
      console.error("Failed to fetch alerts:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAlerts();
  }, []);

  const handleSelect = (alert) => {
    userSelectedRef.current = true;
    setSelected(alert);
  };

  return (
    <div style={styles.container}>
      <h2>🚨 Alerts Center</h2>

      {/* EXPORT */}
      <div style={{ marginBottom: "10px" }}>
        <button onClick={() => exportToCSV(filtered, "alerts.csv")}>
          ⬇ Export CSV
        </button>

        <button onClick={() => exportToJSON(filtered, "alerts.json")}>
          ⬇ Export JSON
        </button>
      </div>

      {/* REFRESH */}
      <div style={{ marginBottom: "10px" }}>
        <button onClick={loadAlerts} style={styles.refreshBtn}>
          🔄 Refresh Alerts
        </button>
      </div>

      {loading && <p>Loading alerts...</p>}

      {!loading && (
        <>
          <Filters alerts={alerts} setFiltered={setFiltered} />

          <AlertTable alerts={filtered} onSelect={handleSelect} />

          {/* 🔥 PASS FULL ALERT */}
          {selected && <AlertDetails alert={selected} />}
        </>
      )}
    </div>
  );
}

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },
  refreshBtn: {
    padding: "6px 12px",
    background: "#2563eb",
    color: "white",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer"
  }
};