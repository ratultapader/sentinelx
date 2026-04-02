import React, { useEffect, useState } from "react";
import KPICards from "../components/KPICards";
import ThreatTrendChart from "../components/ThreatTrendChart";
import FeedbackPanel from "../components/FeedbackPanel";

export default function KPIDashboard() {
  const [stats, setStats] = useState({});
  const [trend, setTrend] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      fetch("http://localhost:9090/api/kpi", {
        headers: { "X-Tenant-ID": "t1" }
      }).then(res => res.json()),

      fetch("http://localhost:9090/api/threat_trend", {
        headers: { "X-Tenant-ID": "t1" }
      }).then(res => res.json())
    ])
      .then(([kpiData, trendData]) => {
        setStats(kpiData);
        setTrend(trendData);
      })
      .catch(err => console.error("KPI load error:", err))
      .finally(() => setLoading(false));
  }, []);

  const handleFeedback = (type) => {
    fetch("http://localhost:9090/api/feedback", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-Tenant-ID": "t1"
      },
      body: JSON.stringify({
        type,
        alert_id: "manual"
      })
    });
  };

  if (loading) return <p style={{ color: "white" }}>Loading...</p>;

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h2>📊 Security KPIs</h2>

      <KPICards stats={stats} />

      <ThreatTrendChart data={trend} />

      <FeedbackPanel onFeedback={handleFeedback} />
    </div>
  );
}