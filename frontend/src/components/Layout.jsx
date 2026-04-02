import React from "react";
import { Link, useLocation } from "react-router-dom";
import TenantBadge from "./TenantBadge";
import TenantSelector from "./TenantSelector"; // ✅ ADD
import ThemeToggle from "./ThemeToggle";

export default function Layout({ children }) {
  const location = useLocation();

  return (
    <div style={styles.container}>
      
      {/* ================= SIDEBAR ================= */}
      <div style={styles.sidebar}>
        <h2 style={styles.logo}>SentinelX</h2>

        <nav>
          <MenuLink to="/" active={location.pathname === "/"}>Dashboard</MenuLink>
          <MenuLink to="/alerts" active={location.pathname === "/alerts"}>Alerts</MenuLink>
          <MenuLink to="/incidents" active={location.pathname === "/incidents"}>Incidents</MenuLink>
          <MenuLink to="/investigation" active={location.pathname === "/investigation"}>Investigation</MenuLink>
          <MenuLink to="/live" active={location.pathname === "/live"}>Live</MenuLink>
          <MenuLink to="/responses" active={location.pathname === "/responses"}>Responses</MenuLink>

          <MenuLink to="/playbooks" active={location.pathname === "/playbooks"}>
  Playbooks
</MenuLink>

<MenuLink to="/attack-map" active={location.pathname === "/attack-map"}>
  🌍 Attack Map
</MenuLink>

          <MenuLink to="/threat-intel" active={location.pathname === "/threat-intel"}>Threat Intel</MenuLink>
          <MenuLink to="/insights" active={location.pathname === "/insights"}>
  Insights
</MenuLink>
          <MenuLink to="/reports" active={location.pathname === "/reports"}>Reports</MenuLink>
          <MenuLink to="/system" active={location.pathname === "/system"}>System</MenuLink>

          <MenuLink to="/admin" active={location.pathname === "/admin"}>
  🛡 Admin
</MenuLink>

          <MenuLink to="/kpi" active={location.pathname === "/kpi"}>
  📊 KPI Dashboard
</MenuLink>

          <MenuLink to="/rules" active={location.pathname === "/rules"}>
  Rules
</MenuLink>
        </nav>
      </div>

      {/* ================= MAIN ================= */}
      <div style={styles.main}>
        
        {/* 🔥 HEADER */}
        <div style={styles.header}>
  <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
    <ThemeToggle />   {/* 🔥 ADD THIS */}
    <TenantSelector />
  </div>

  <TenantBadge />
</div>

        {/* 🔥 CONTENT */}
        <div style={styles.content} className="fade-in">
          {children}
        </div>

      </div>

    </div>
  );
}

/* ================= MENU LINK ================= */
function MenuLink({ to, children, active }) {
  return (
    <Link
      to={to}
      style={{
        ...styles.link,
        background: active ? "#1e293b" : "transparent",
        borderLeft: active ? "3px solid #3b82f6" : "3px solid transparent"
      }}
    >
      {children}
    </Link>
  );
}

/* ================= STYLES ================= */
const styles = {
  container: {
    display: "flex",
    height: "100vh",
    overflow: "hidden"
  },
  sidebar: {
    width: "220px",
    background: "#0f172a",
    color: "#e2e8f0",
    padding: "20px",
    borderRight: "1px solid #1e293b"
  },
  logo: {
    marginBottom: "20px",
    fontWeight: "600"
  },
  main: {
    flex: 1,
    display: "flex",
    flexDirection: "column"
  },
  header: {
    height: "60px",
    background: "#1e293b",
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between", // ✅ FIXED
    padding: "0 20px",
    borderBottom: "1px solid #334155"
  },
  content: {
    flex: 1,
    padding: "20px",
    background: "#020617",
    overflowY: "auto"
  },
  link: {
    display: "block",
    margin: "6px 0",
    padding: "10px",
    borderRadius: "6px",
    color: "#e2e8f0",
    textDecoration: "none",
    transition: "0.2s"
  }
};