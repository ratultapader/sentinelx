import React from "react";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

function AlertsTimeline({ data }) {
  if (!data || data.length === 0) return null;

  return (
    <div>
      <h3>Alerts Timeline</h3>
      <ResponsiveContainer width="100%" height={250}>
        <LineChart data={data}>
          <XAxis dataKey="time" stroke="#94a3b8" />
          <YAxis stroke="#94a3b8" />
          <Tooltip />
          <Line type="monotone" dataKey="avg_anomaly_score" stroke="#22c55e" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

export default AlertsTimeline;