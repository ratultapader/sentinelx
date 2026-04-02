import React, { useState, useEffect } from "react";

export default function AlertDetails({ alert }) {
  const [status, setStatus] = useState("NEW");
  const [note, setNote] = useState("");
  const [notes, setNotes] = useState([]);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (alert) {
      setStatus(alert.status || "NEW");
      setNotes(alert.notes || []);
      setNote("");
    }
  }, [alert]);

  if (!alert) {
    return <div style={{ color: "#94a3b8" }}>No alert selected</div>;
  }

  const statusColor = {
    NEW: "#dc2626",
    INVESTIGATING: "#f97316",
    RESOLVED: "#22c55e"
  };

  const saveUpdate = async () => {
    try {
      const alertId = alert?.id || alert?._id;

      if (!alertId) {
        console.error("❌ Missing alert ID", alert);
        return;
      }

      setSaving(true);

      const res = await fetch(
        `http://localhost:9090/api/alerts/${alertId}`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "X-Tenant-ID": "t1"
          },
          body: JSON.stringify({
            status,
            note
          })
        }
      );

      if (!res.ok) throw new Error("Update failed");

      if (note.trim()) {
        setNotes(prev => [...prev, note]);
        setNote("");
      }

      window.alert("✅ Alert updated successfully");

    } catch (err) {
      console.error("Update error:", err);
      window.alert("❌ Failed to update alert");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div style={styles.card}>
      <h3>🚨 Alert Details</h3>

      <p><strong>ID:</strong> {alert.id}</p>
      <p><strong>Type:</strong> {alert.type || alert.event_type}</p>
      <p><strong>Source IP:</strong> {alert.source_ip}</p>

      {/* 🔥 NEW: THREAT SCORE */}
      <h3 style={{ color: "#ef4444", marginTop: "10px" }}>
        🔥 Threat Score: {alert.threat_score}
      </h3>

      {/* 🔥 NEW: SCORING BARS */}
      {renderBar("Anomaly Detection", alert.anomaly_score)}
      {renderBar("Signature Match", alert.signature_score)}
      {renderBar("IP Reputation", alert.ip_reputation)}
      {renderBar("Behavior Analysis", alert.behavior_score)}

      {/* 🔥 STATUS */}
      <div style={{ marginTop: "15px" }}>
        <strong>Status: </strong>

        <select
          value={status}
          onChange={e => setStatus(e.target.value)}
          style={{
            marginLeft: "10px",
            padding: "5px",
            background: "#020617",
            color: statusColor[status],
            border: "1px solid #334155"
          }}
        >
          <option value="NEW">NEW</option>
          <option value="INVESTIGATING">INVESTIGATING</option>
          <option value="RESOLVED">RESOLVED</option>
        </select>
      </div>

      {/* NOTES INPUT */}
      <div style={{ marginTop: "15px" }}>
        <strong>Notes:</strong>

        <textarea
          value={note}
          onChange={e => setNote(e.target.value)}
          placeholder="Add investigation notes..."
          style={styles.textarea}
        />
      </div>

      {/* SAVE */}
      <button
        onClick={saveUpdate}
        disabled={saving}
        style={styles.button}
      >
        {saving ? "Saving..." : "💾 Save"}
      </button>

      {/* PREVIOUS NOTES */}
      {notes.length > 0 && (
        <div style={{ marginTop: "15px" }}>
          <strong>Previous Notes:</strong>

          {notes.map((n, i) => (
            <div key={i} style={styles.note}>
              {n}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

//////////////////////////////////////////////////////
// 🔥 SCORING BAR COMPONENT
//////////////////////////////////////////////////////

function renderBar(label, value) {
  return (
    <div style={{ marginTop: "10px" }}>
      <p>{label} ({value || 0})</p>

      <div style={styles.barBg}>
        <div
          style={{
            ...styles.barFill,
            width: `${(value || 0) * 100}%`
          }}
        />
      </div>
    </div>
  );
}

//////////////////////////////////////////////////////

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px",
    color: "#e2e8f0",
    marginTop: "10px"
  },
  textarea: {
    width: "100%",
    height: "80px",
    marginTop: "5px",
    background: "#020617",
    color: "white",
    border: "1px solid #334155",
    padding: "8px"
  },
  button: {
    marginTop: "10px",
    padding: "8px 12px",
    background: "#2563eb",
    color: "white",
    border: "none",
    borderRadius: "5px",
    cursor: "pointer"
  },
  note: {
    background: "#334155",
    marginTop: "5px",
    padding: "6px",
    borderRadius: "5px"
  },

  // 🔥 NEW
  barBg: {
    height: "8px",
    background: "#374151",
    borderRadius: "5px"
  },
  barFill: {
    height: "100%",
    background: "#3b82f6",
    borderRadius: "5px",
    transition: "width 0.4s ease"
  }
};