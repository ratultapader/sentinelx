import React, { useEffect, useState } from "react";

export default function TenantSelector() {
  const [tenant, setTenant] = useState("t1");

  useEffect(() => {
    const saved = localStorage.getItem("tenant_id");
    if (saved) setTenant(saved);
  }, []);

  const changeTenant = (t) => {
    setTenant(t);
    localStorage.setItem("tenant_id", t);
    window.location.reload(); // simple + reliable
  };

  return (
    <div style={styles.container}>
      <strong>Tenant:</strong>

      <select
        value={tenant}
        onChange={(e) => changeTenant(e.target.value)}
        style={styles.select}
      >
        <option value="t1">Tenant 1</option>
        <option value="t2">Tenant 2</option>
        <option value="t3">Tenant 3</option>
      </select>
    </div>
  );
}

const styles = {
  container: {
    display: "flex",
    alignItems: "center",
    gap: "10px",
    marginBottom: "10px",
    color: "white"
  },
  select: {
    background: "#1e293b",
    color: "white",
    padding: "5px",
    borderRadius: "6px",
    border: "none"
  }
};