import React from "react";

export default function TenantBadge() {
  const tenant = localStorage.getItem("tenant_id") || "t1";

  return (
    <div style={{
      background: "#334155",
      padding: "5px 10px",
      borderRadius: "8px",
      color: "white"
    }}>
      Active Tenant: {tenant}
    </div>
  );
}