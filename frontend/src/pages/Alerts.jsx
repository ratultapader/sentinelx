import React, { useEffect, useState, useRef } from "react";
import { fetchRecentAlerts } from "../services/dashboardApi";
import Filters from "../components/Filters";
import AlertTable from "../components/AlertTable";
import AlertDetails from "../components/AlertDetails";

export default function Alerts() {
  const [alerts, setAlerts] = useState([]);
  const [filtered, setFiltered] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);

  // 🔥 track user manual selection (prevents override)
  const userSelectedRef = useRef(false);

  // ✅ LOAD ALERTS
  const loadAlerts = async () => {
    try {
      setLoading(true);

      const data = await fetchRecentAlerts();
      console.log("RAW ALERTS:", data);

      let items = Array.isArray(data) ? data : data?.items || [];

      // 🔥 NORMALIZE DATA (CRITICAL)
      const normalized = items.map((item, index) => ({
        id: item.id || item._id || `alert-${index}`,
        type: item.event_type || item.type || "Unknown",
        source_ip: item.source_ip || item.ip || "N/A",
        severity: item.severity || "low",
        timestamp: item.timestamp || new Date().toISOString(),
        ...item
      }));

      setAlerts(normalized);
      setFiltered(normalized);

      // ✅ only auto select FIRST TIME (not after user clicks)
      if (normalized.length > 0 && !userSelectedRef.current) {
        setSelected(normalized[0]);
      }

    } catch (err) {
      console.error("Failed to fetch alerts:", err);
    } finally {
      setLoading(false);
    }
  };

  // ✅ LOAD ONCE (NO AUTO REFRESH BUG)
  useEffect(() => {
    loadAlerts();
  }, []);

  // ✅ HANDLE USER CLICK (VERY IMPORTANT)
  const handleSelect = (alert) => {
    userSelectedRef.current = true;
    setSelected(alert);
  };

  return (
    <div style={styles.container}>
      <h2>🚨 Alerts Center</h2>

      {/* 🔄 MANUAL REFRESH BUTTON */}
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