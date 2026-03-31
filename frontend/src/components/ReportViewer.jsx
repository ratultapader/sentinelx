import React from "react";

export default function ReportViewer({ report }) {
  const timeline = report.attack_chain || report.timeline || [];
  const severity = report.severity || report.risk_level || "unknown";

  const severityColor = (sev) => {
    const s = (sev || "").toLowerCase();
    if (s === "critical") return "#ef4444";
    if (s === "high") return "#f97316";
    if (s === "medium") return "#eab308";
    if (s === "low") return "#22c55e";
    return "#94a3b8";
  };

  // ✅ FIXED PDF DOWNLOAD
  const downloadPDF = async () => {
    try {
      console.log("Downloading PDF for:", report.incident_id);

      const res = await fetch(
        `http://localhost:9090/api/reports/${report.incident_id}/pdf`,
        {
          headers: {
            "X-Tenant-ID": "t1",
          },
        }
      );

      if (!res.ok) throw new Error("PDF failed");

      const blob = await res.blob();
      const url = window.URL.createObjectURL(blob);
      window.open(url);
    } catch (err) {
      console.error("PDF ERROR:", err);
      alert("Failed to download PDF");
    }
  };

  return (
    <div style={styles.container}>
      {/* HEADER */}
      <div style={styles.header}>
        <div>
          <h2>📄 Incident Report</h2>
          <p style={{ opacity: 0.6 }}>ID: {report.incident_id}</p>
        </div>

        <div
          style={{
            ...styles.badge,
            background: severityColor(severity),
          }}
        >
          {severity.toUpperCase()}
        </div>
      </div>

      {/* GRID */}
      <div style={styles.grid}>
        <Card title="Executive Summary">
          {report.executive_summary || "No summary available"}
        </Card>

        <Card title="Attack Type">
          {report.attack_type || "Unknown"}
        </Card>

        <Card title="MITRE">
          <div>
            <b>Tactic:</b> {report.mitre_tactic || "-"}
          </div>
          <div>
            <b>Technique:</b> {report.mitre_technique || "-"}
          </div>
        </Card>

        <Card title="Source">
          {report.source_ip || "N/A"}
        </Card>
      </div>

      {/* TIMELINE */}
      <div style={styles.section}>
        <h3>🕒 Attack Timeline</h3>

        {timeline.length === 0 ? (
          <p style={{ opacity: 0.6 }}>No timeline available</p>
        ) : (
          timeline.map((t, i) => (
            <div key={i} style={styles.timelineItem}>
              <div style={styles.timelineDot} />
              <div>
                <div style={{ fontWeight: "bold" }}>
                  {t.timestamp || "time"}
                </div>
                <div>{t.event_type || t.stage}</div>
                <div style={{ opacity: 0.7 }}>
                  {t.summary || "Suspicious activity"}
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      {/* RECOMMENDATIONS */}
      <div style={styles.section}>
        <h3>🛡 Recommendations</h3>
        <ul>
          {(report.recommended_remediation || [
            "Block malicious IP",
            "Enable rate limiting",
            "Monitor suspicious activity",
          ]).map((r, i) => (
            <li key={i}>{r}</li>
          ))}
        </ul>
      </div>

      {/* BUTTON */}
      <div style={{ textAlign: "center", marginTop: "30px" }}>
        <button style={styles.button} onClick={downloadPDF}>
          ⬇ Download PDF
        </button>
      </div>
    </div>
  );
}

/* ================= STYLES ================= */

const styles = {
  container: {
    marginTop: "20px",
    background: "#0f172a",
    padding: "30px",
    borderRadius: "12px",
    color: "#e2e8f0",
    maxWidth: "1000px",
    marginInline: "auto",
    boxShadow: "0 10px 30px rgba(0,0,0,0.4)",
  },

  header: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: "25px",
  },

  badge: {
    padding: "8px 14px",
    borderRadius: "20px",
    fontWeight: "bold",
    color: "#fff",
  },

  grid: {
    display: "grid",
    gridTemplateColumns: "repeat(2, 1fr)",
    gap: "15px",
    marginBottom: "25px",
  },

  section: {
    marginBottom: "25px",
  },

  timelineItem: {
    display: "flex",
    gap: "12px",
    marginBottom: "12px",
    background: "#1e293b",
    padding: "12px",
    borderRadius: "8px",
  },

  timelineDot: {
    width: "10px",
    height: "10px",
    background: "#3b82f6",
    borderRadius: "50%",
    marginTop: "6px",
  },

  button: {
    padding: "12px 22px",
    borderRadius: "8px",
    border: "none",
    background: "#3b82f6",
    color: "white",
    cursor: "pointer",
    fontWeight: "bold",
  },
};

/* CARD */
function Card({ title, children }) {
  return (
    <div
      style={{
        background: "#1e293b",
        padding: "15px",
        borderRadius: "10px",
      }}
    >
      <h4 style={{ marginBottom: "10px" }}>{title}</h4>
      <div>{children}</div>
    </div>
  );
}