import React from "react";

export default function LiveFeed({ events }) {
  if (!events || events.length === 0) {
    return <p style={{ color: "#94a3b8" }}>No live events</p>;
  }

  return (
    <div>
      {events.map((e, i) => (
        <div key={i} style={getStyle(e)}>
          
          {/* MAIN LINE */}
          <div>
            <strong>{e.type || "event"}</strong> —{" "}
            {e.description || e.target || "No data"}
          </div>

          {/* EXTRA INFO (SOC STYLE 🔥) */}
          <div style={{ fontSize: "12px", opacity: 0.8 }}>
            {e.source_ip && <span>IP: {e.source_ip} | </span>}
            {e.threat_score !== undefined && (
              <span>Score: {e.threat_score}</span>
            )}
          </div>

        </div>
      ))}
    </div>
  );
}

function getStyle(e) {
  if (e.severity === "critical") {
    return {
      background: "#dc2626",
      color: "white",
      padding: "12px",
      marginBottom: "8px",
      borderRadius: "8px",
      animation: "pulse 1s infinite"
    };
  }

  if (e.severity === "high") {
    return {
      background: "#7c2d12",
      color: "#fff",
      padding: "12px",
      marginBottom: "8px",
      borderRadius: "8px"
    };
  }

  return {
    background: "#1e293b",
    color: "#e2e8f0",
    padding: "12px",
    marginBottom: "8px",
    borderRadius: "8px"
  };
}