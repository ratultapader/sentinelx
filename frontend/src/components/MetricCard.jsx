import { useNavigate } from "react-router-dom";

export default function MetricCard({ title, value }) {
  const navigate = useNavigate();

  return (
    <div
      className="card"
      onClick={() => navigate("/alerts")} // ✅ navigate to Alerts Center
      style={{
        cursor: "pointer",
        transition: "all 0.2s ease",
      }}
      onMouseEnter={(e) => {
        e.currentTarget.style.transform = "scale(1.03)";
      }}
      onMouseLeave={(e) => {
        e.currentTarget.style.transform = "scale(1)";
      }}
    >
      <h3>{title}</h3>
      <h2>{value}</h2>
    </div>
  );
}