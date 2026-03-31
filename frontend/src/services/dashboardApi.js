import api from "./api";

// ===============================
// DASHBOARD
// ===============================

export const fetchDashboard = async () => {
  const res = await api.get("/api/dashboard/anomalies");
  return res.data;
};

// ===============================
// ALERTS
// ===============================

export const fetchRecentAlerts = async () => {
  const res = await api.get("/api/alerts/recent");
  return res.data;
};

// ===============================
// INCIDENTS
// ===============================

export const fetchIncidents = async () => {
  const res = await api.get("/api/incidents");
  return res.data;
};

export const fetchIncidentById = async (id) => {
  if (!id) throw new Error("Invalid incident ID");

  const res = await api.get(`/api/incidents/${id}`);
  return res.data;
};