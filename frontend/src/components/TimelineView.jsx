import React, { useEffect, useState } from "react";
import api from "../services/api";

export default function TimelineView({ ip }) {
  const [events, setEvents] = useState([]);

  useEffect(() => {
    if (!ip) return;

    api.get(`/api/timeline/${ip}`, {
      headers: { "X-Tenant-ID": "t1" }
    })
    .then(res => {
      const data = res.data;

      const merged = [
        ...(data.alerts || []),
        ...(data.events || []),
        ...(data.actions || [])
      ];

      merged.sort((a, b) =>
        new Date(a.timestamp) - new Date(b.timestamp)
      );

      setEvents(merged);
    })
    .catch(() => setEvents([]));

  }, [ip]);

  return (
    <div>
      {events.length > 0 ? (
        events.map((e, i) => (
          <div key={i} style={styles.event}>
            
            <div style={styles.time}>
              {formatTime(e.timestamp)}
            </div>

            <div style={styles.content}>
              <strong>{e.type || "event"}</strong>
              <p>{e.destination || e.source_ip}</p>
            </div>

          </div>
        ))
      ) : (
        <p>No timeline data</p>
      )}
    </div>
  );
}

function formatTime(ts) {
  return ts ? new Date(ts).toLocaleTimeString() : "";
}

const styles = {
  event: {
    display: "flex",
    marginBottom: "10px",
    background: "#1e293b",
    padding: "10px",
    borderRadius: "8px"
  },
  time: {
    width: "80px",
    fontSize: "12px",
    color: "#94a3b8"
  },
  content: {
    marginLeft: "10px",
    color: "#e2e8f0"
  }
};