export function isAdmin() {
  return localStorage.getItem("role") === "admin";
}