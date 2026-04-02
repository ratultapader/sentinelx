import React from "react";
import AuditLogs from "../components/AuditLogs";
import PerformancePanel from "../components/PerformancePanel";
import APIExplorer from "../components/APIExplorer";
import { isAdmin } from "../utils/auth";

export default function Admin() {
  // 🔒 ROLE PROTECTION
  if (!isAdmin()) {
    return <p style={{ color: "red" }}>Access Denied</p>;
  }

  return (
    <div style={styles.container}>
      <h2 style={styles.title}>🛡 Admin Panel</h2>

      <AuditLogs />
      <PerformancePanel />
      <APIExplorer />
    </div>
  );
}

const styles = {
  container: {
    padding: "20px",
    color: "white"
  },
  title: {
    marginBottom: "20px"
  }
};