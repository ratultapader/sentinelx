import axios from "axios";

const api = axios.create({
 baseURL: "http://localhost:9090",
  timeout: 10000,
});

// 🔥 Tenant interceptor
api.interceptors.request.use((config) => {
  const tenant = localStorage.getItem("tenant_id") || "t1";

  config.headers["X-Tenant-ID"] = tenant;
  config.headers["Content-Type"] = "application/json";

  return config;
});

// 🔥 Error handler
api.interceptors.response.use(
  (res) => res,
  (err) => {
    console.error("API ERROR:", err?.response || err.message);
    alert("API Error");
    return Promise.reject(err);
  }
);

// ==============================
// 🔥 EXISTING HELPERS (NO CHANGE)
// ==============================

export const get = async (url) => {
  const res = await api.get(url);
  return res.data;
};

export const post = async (url, body) => {
  const res = await api.post(url, body);
  return res.data;
};

export const postNoBody = async (url) => {
  const res = await api.post(url);
  return res.data;
};

// ==============================
// 🚀 DAY 64 ADDITION (SAFE)
// ==============================

// 👉 Use this ONLY when API returns { count, items }
export const getItems = async (url) => {
  const data = await get(url);

  if (data?.items) return data.items; // incidents case
  return data; // alerts case (array)
};

export default api;