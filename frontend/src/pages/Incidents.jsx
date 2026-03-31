import React, { useEffect, useState, useRef } from "react";
import { fetchIncidents } from "../services/dashboardApi";
import IncidentTable from "../components/IncidentTable";
import IncidentDetails from "../components/IncidentDetails";

export default function Incidents() {
  const [incidents, setIncidents] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // 🔥 keep selected id stable across refresh
  const selectedIdRef = useRef(null);

  const loadIncidents = async () => {
    try {
      setError(null);

      const data = await fetchIncidents();
      console.log("INCIDENT API RESPONSE:", data);

      let items = Array.isArray(data?.items) ? data.items : [];

      // ✅ sort latest first
      items.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));

      setIncidents(items);

      // 🔥 RESTORE SELECTION (PRODUCTION FIX)
      if (selectedIdRef.current) {
        const found = items.find(i => i.id === selectedIdRef.current);
        if (found) {
          setSelected(found);
          return;
        }
      }

      // ✅ fallback: select first only if nothing selected
      if (items.length > 0 && !selectedIdRef.current) {
        setSelected(items[0]);
        selectedIdRef.current = items[0].id;
      }

      setLoading(false);

    } catch (err) {
      console.error("Failed to fetch incidents:", err);
      setError("Failed to load incidents");
      setLoading(false);
    }
  };

  // 🔥 track user selection
  const handleSelect = (incident) => {
    setSelected(incident);
    selectedIdRef.current = incident.id;
  };

  useEffect(() => {
    loadIncidents();

    const interval = setInterval(loadIncidents, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h1>📂 Incidents Center</h1>

      {loading && <p>Loading incidents...</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}

      {!loading && !error && (
        <>
          <IncidentTable
            incidents={incidents}
            onSelect={handleSelect}   // 🔥 FIXED
          />

          {selected && (
            <IncidentDetails incident={selected} />
          )}
        </>
      )}
    </div>
  );
}