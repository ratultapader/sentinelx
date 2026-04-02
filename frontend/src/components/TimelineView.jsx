import React, { useEffect, useState } from "react";

export default function TimelineView({ ip }) {
  const [events, setEvents] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!ip) return;

    setLoading(true);

    fetch(`http://localhost:9090/api/timeline/${ip}`, {
      headers: { "X-Tenant-ID": "t1" }
    })
      .then(res => res.json())
      .then(data => {
        console.log("TIMELINE:", data);

        // ✅ FIX: handle new array format
        if (Array.isArray(data)) {
          setEvents(data);
        } else {
          setEvents([]);
        }
      })
      .catch(err => {
        console.error("Timeline error:", err);
        setEvents([]);
      })
      .finally(() => setLoading(false));

  }, [ip]);

  if (loading) return <p style={styles.muted}>Loading timeline...</p>;

  if (!events.length) {
    return <p style={styles.muted}>No timeline data</p>;
  }

  return (
    <div style={styles.container}>
      {events.map((e, i) => (
        <div key={i} style={styles.event}>
          <span style={styles.time}>{formatTime(e.timestamp)}</span>
          <span style={styles.arrow}>→</span>
          <span style={styles.type}>{e.type}</span>
        </div>
      ))}
    </div>
  );
}

// 🔥 FORMAT TIME
function formatTime(ts) {
  return new Date(ts).toLocaleTimeString();
}

const styles = {
  container: {
    display: "flex",
    flexDirection: "column",
    gap: "6px"
  },
  event: {
    background: "#1e293b",
    padding: "8px",
    borderRadius: "6px",
    display: "flex",
    alignItems: "center"
  },
  time: {
    color: "#38bdf8",
    marginRight: "8px",
    fontSize: "12px"
  },
  arrow: {
    marginRight: "8px"
  },
  type: {
    color: "#e2e8f0"
  },
  muted: {
    color: "#94a3b8"
  }
};