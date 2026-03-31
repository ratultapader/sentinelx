import React, { useEffect, useState } from "react";
import api from "../services/api";
import ReportViewer from "../components/ReportViewer";

export default function Reports() {
  const [reports, setReports] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(false);

  // ✅ Fetch report list
  useEffect(() => {
    api.get("/api/reports", {
      headers: { "X-Tenant-ID": "t1" }
    })
      .then(res => setReports(res.data.items || []))
      .catch(() => setReports([]));
  }, []);

  // ✅ Fetch FULL report when clicking View
  const handleView = async (incidentId) => {
    try {
      setLoading(true);
      const res = await api.get(`/api/reports/${incidentId}`, {
        headers: { "X-Tenant-ID": "t1" }
      });
      setSelected(res.data);
    } catch (err) {
      console.error("Failed to load report", err);
      alert("Failed to load report");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={styles.container}>
      <h2 style={styles.title}>📊 Reports</h2>

      {/* TABLE */}
      <table style={styles.table}>
        <thead>
          <tr>
            <th>Incident ID</th>
            <th>Severity</th>
            <th>Created</th>
            <th>Action</th>
          </tr>
        </thead>

        <tbody>
          {reports.map(r => (
            <tr key={r.id} style={styles.row}>
              <td>{r.incident_id}</td>
              <td style={{ color: getSeverityColor(r.severity) }}>
                {r.severity}
              </td>
              <td>{new Date(r.created_at).toLocaleString()}</td>
              <td>
                <button
                  style={styles.viewBtn}
                  onClick={() => handleView(r.incident_id)}
                >
                  View
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {/* LOADING */}
      {loading && (
        <div style={styles.loading}>
          Loading report...
        </div>
      )}

      {/* REPORT VIEW */}
      {selected && !loading && (
        <ReportViewer report={selected} />
      )}
    </div>
  );
}

/* ================= STYLES ================= */

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },

  title: {
    marginBottom: "15px"
  },

  table: {
    width: "100%",
    borderCollapse: "collapse",
    background: "#0f172a",
    borderRadius: "8px",
    overflow: "hidden"
  },

  row: {
    borderBottom: "1px solid #1e293b"
  },

  viewBtn: {
    padding: "6px 12px",
    background: "#3b82f6",
    border: "none",
    borderRadius: "6px",
    color: "white",
    cursor: "pointer"
  },

  loading: {
    marginTop: "20px",
    textAlign: "center",
    opacity: 0.7
  }
};

/* ================= HELPERS ================= */

function getSeverityColor(sev) {
  const s = (sev || "").toLowerCase();
  if (s === "critical") return "#ef4444";
  if (s === "high") return "#f97316";
  if (s === "medium") return "#eab308";
  if (s === "low") return "#22c55e";
  return "#94a3b8";
}