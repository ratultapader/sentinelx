import React, { useState } from "react";

export default function Filters({ alerts, setFiltered }) {
  const [severity, setSeverity] = useState("");
  const [search, setSearch] = useState("");

  const applyFilters = () => {
    let result = alerts;

    if (severity) {
      result = result.filter(a => a.severity === severity);
    }

    if (search) {
      result = result.filter(a =>
        a.source_ip.toLowerCase().includes(search.toLowerCase())
      );
    }

    setFiltered(result);
  };

  return (
    <div style={{ marginBottom: "15px" }}>
      <select onChange={e => setSeverity(e.target.value)}>
        <option value="">All Severity</option>
        <option value="critical">Critical</option>
        <option value="high">High</option>
        <option value="medium">Medium</option>
        <option value="low">Low</option>
      </select>

      <input
        placeholder="Search IP..."
        onChange={e => setSearch(e.target.value)}
        style={{ marginLeft: "10px" }}
      />

      <button onClick={applyFilters} style={{ marginLeft: "10px" }}>
        Apply
      </button>
    </div>
  );
}