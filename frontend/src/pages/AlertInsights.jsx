import React, { useEffect, useState } from "react";
import api from "../services/api";
import PriorityCard from "../components/PriorityCard";
import AttackStory from "../components/AttackStory";

export default function AlertInsights() {
  const [data, setData] = useState(null);
  const [error, setError] = useState(null);

  const loadData = async () => {
    try {
      const res = await api.get("/api/alert_insights"); // ✅ tenant auto
      setData(res.data);
    } catch (err) {
      console.error(err);
      setError("Failed to load insights");
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  if (error) return <p>{error}</p>;
  if (!data) return <p>Loading insights...</p>;

  return (
    <div style={{ padding: "20px", color: "white" }}>
      <h2>🧠 Alert Insights</h2>

      <PriorityCard
        score={data.priority_score}
        reasons={data.reasons}
      />

      {/* 🔥 FIXED HERE */}
      <AttackStory steps={data.story_steps || []} />
    </div>
  );
}