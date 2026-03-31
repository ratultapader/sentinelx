import React, { useEffect, useState } from "react";
import HealthCards from "../components/HealthCards";
import MetricsPanel from "../components/MetricsPanel";
import api from "../services/api"; // ✅ use axios instance

export default function System() {
  const [health, setHealth] = useState({});
  const [metrics, setMetrics] = useState({});
  const [loading, setLoading] = useState(true);

  const loadData = async () => {
    try {
      const [healthRes, metricsRes] = await Promise.all([
        api.get("/health"),   // ✅ tenant header auto added
        api.get("/metrics")
      ]);

      const services = healthRes.data.services || {};

      // ✅ map backend → frontend
      setHealth({
        database: services.database === "up" ? "ok" : "down",
        elasticsearch: services.elasticsearch === "up" ? "ok" : "down",
        neo4j: services.neo4j === "up" ? "ok" : "down"
      });

      setMetrics(metricsRes.data || {});
      setLoading(false);

    } catch (err) {
      console.error("System load error:", err);
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();

    const interval = setInterval(loadData, 5000); // 🔄 auto refresh
    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h2>🖥️ System Health</h2>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <HealthCards health={health} />
      )}

      <h2 style={{ marginTop: "30px" }}>📊 Metrics</h2>

      <MetricsPanel metrics={metrics} />
    </div>
  );
}