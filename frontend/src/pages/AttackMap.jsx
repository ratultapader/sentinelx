import React, { useEffect, useState } from "react";
import {
  MapContainer,
  TileLayer,
  CircleMarker,
  Popup
} from "react-leaflet";

export default function AttackMap() {
  const [attacks, setAttacks] = useState([]);

  // ===============================
  // FETCH + MERGE (NO DUPLICATES)
  // ===============================
  const loadAttacks = () => {
    fetch("http://localhost:9090/api/attack_map", {
      headers: { "X-Tenant-ID": "t1" }
    })
      .then(res => res.json())
      .then(newData => {
        setAttacks(prev => {
          const existing = new Set(
            prev.map(a => `${a.ip}-${a.lat}-${a.lng}`)
          );

          const newPoints = newData.filter(
            a => !existing.has(`${a.ip}-${a.lat}-${a.lng}`)
          );

          return [...prev, ...newPoints];
        });
      })
      .catch(err => console.error(err));
  };

  useEffect(() => {
    loadAttacks();
  }, []);

  return (
    <div style={{ height: "100%", padding: "10px" }}>

      {/* HEADER */}
      <div style={styles.header}>
        <h2>🌍 Attack Map ({attacks.length})</h2>

        <button onClick={loadAttacks} style={styles.btn}>
          🔄 Refresh
        </button>
      </div>

      {/* MAP */}
      <MapContainer
        center={[20, 0]}
        zoom={2}
        scrollWheelZoom={true}
        style={styles.map}
      >
        <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />

        {/* 🔴 ATTACK POINTS */}
        {attacks.map((a, i) => (
          <CircleMarker
            key={`${a.ip}-${a.lat}-${a.lng}-${i}`}
            center={[a.lat, a.lng]}
            radius={6}
            pathOptions={{
              color: "red",
              fillColor: "red",
              fillOpacity: 0.8,
              weight: 2
            }}
          >
            {/* 🔥 ADVANCED POPUP */}
            <Popup>
              <div style={styles.popup}>
                <h4 style={{ margin: "0 0 5px 0" }}>
                  🌍 {a.country || "Unknown"}
                </h4>

                <p><strong>IP:</strong> {a.ip}</p>

                {a.target && (
                  <p><strong>Target:</strong> {a.target}</p>
                )}

                {a.threat_score && (
                  <p>
                    <strong>Threat Score:</strong>{" "}
                    {a.threat_score.toFixed(2)}
                  </p>
                )}

                {a.timestamp && (
                  <p>
                    <strong>Time:</strong>{" "}
                    {new Date(a.timestamp).toLocaleString()}
                  </p>
                )}
              </div>
            </Popup>

          </CircleMarker>
        ))}
      </MapContainer>
    </div>
  );
}

// ===============================
// STYLES
// ===============================
const styles = {
  header: {
    display: "flex",
    justifyContent: "space-between",
    marginBottom: "10px",
    alignItems: "center"
  },
  btn: {
    padding: "6px 12px",
    border: "none",
    borderRadius: "6px",
    background: "#1e293b",
    color: "#fff",
    cursor: "pointer"
  },
  map: {
    height: "90vh",
    borderRadius: "10px"
  },
  popup: {
    fontSize: "13px",
    lineHeight: "1.5"
  }
};