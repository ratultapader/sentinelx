import { BrowserRouter, Routes, Route } from "react-router-dom";

import Dashboard from "./pages/Dashboard";
import Alerts from "./pages/Alerts";
import Incidents from "./pages/Incidents";
import Investigation from "./pages/Investigation";
import Live from "./pages/Live";
import Responses from "./pages/Responses";
import Reports from "./pages/Reports";
import ThreatIntel from "./pages/ThreatIntel";
import System from "./pages/System";
import AlertInsights from "./pages/AlertInsights";
import Rules from "./pages/Rules";
import Playbooks from "./pages/Playbooks";
import AttackMap from "./pages/AttackMap";
// import Investigation from "./pages/Investigation";
import Admin from "./pages/Admin";

import KPIDashboard from "./pages/KPIDashboard";


// ✅ LAYOUT (IMPORTANT)
import Layout from "./components/Layout";

export default function App() {
  return (
    <BrowserRouter>

      {/* ✅ FULL APP WRAPPED IN LAYOUT */}
      <Layout>

        <Routes>

          {/* 🏠 Dashboard */}
          <Route path="/" element={<Dashboard />} />

          {/* 🚨 Alerts */}
          <Route path="/alerts" element={<Alerts />} />

          {/* 📂 Incidents */}
          <Route path="/incidents" element={<Incidents />} />

          {/* 🔍 Investigation */}
          <Route path="/investigation" element={<Investigation />} />

          {/* 🔴 Live Monitoring */}
          <Route path="/live" element={<Live />} />

          {/* ⚡ Responses */}
          <Route path="/responses" element={<Responses />} />

          {/* 📄 Reports */}
          <Route path="/reports" element={<Reports />} />

          {/* 🧠 Threat Intel */}
          <Route path="/threat-intel" element={<ThreatIntel />} />

          {/* 🖥️ System */}
          <Route path="/system" element={<System />} />

          <Route path="/insights" element={<AlertInsights />} />

          <Route path="/rules" element={<Rules />} />

          <Route path="/playbooks" element={<Playbooks />} />

          <Route path="/attack-map" element={<AttackMap />} />

          <Route path="/investigation" element={<Investigation />} />

          <Route path="/kpi" element={<KPIDashboard />} />

          <Route path="/admin" element={<Admin />} />

        </Routes>

      </Layout>

    </BrowserRouter>
  );
}