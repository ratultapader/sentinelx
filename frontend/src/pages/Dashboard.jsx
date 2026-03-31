import { useEffect, useState } from "react";
import "../styles/main.css";
import MetricCard from "../components/MetricCard";
import { fetchDashboard, fetchRecentAlerts } from "../services/dashboardApi";
import SeverityChart from "../components/dashboard/SeverityChart";
import AlertsTimeline from "../components/dashboard/AlertsTimeline";
import RecentAlerts from "../components/dashboard/RecentAlerts";
import ThreatPieChart from "../components/dashboard/ThreatPieChart";
import "../styles/dashboard.css";

export default function Dashboard() {
  const [data, setData] = useState(null);
  const [alerts, setAlerts] = useState([]);

  useEffect(() => {
    fetchDashboard().then(setData);
    fetchRecentAlerts().then(setAlerts);
  }, []);

  if (!data) return <div className="container">Loading...</div>;

  return (
    <div className="container">
      <h1>SentinelX SOC Dashboard</h1>

      {/* METRICS */}
      <div className="dashboard-grid">
        <MetricCard title="Total Alerts" value={data.summary.total_alerts} />
        <MetricCard title="Critical" value={data.summary.critical_alerts} />
        <MetricCard title="High" value={data.summary.high_alerts} />
      </div>

      {/* CHARTS */}
      <div className="chart-grid">
        <ThreatPieChart data={data.threat_distribution} />
        <SeverityChart summary={data.summary} />
        <AlertsTimeline data={data.anomaly_trend} />
      </div>

      {/* RECENT ALERTS */}
      <RecentAlerts alerts={alerts} summary={data.summary} />
    </div>
  );
}