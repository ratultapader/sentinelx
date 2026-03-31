import React from "react";
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

function SeverityChart({ summary }) {
  if (!summary) return null;

  const data = [
    { name: "Critical", value: summary.critical_alerts },
    { name: "High", value: summary.high_alerts },
    { name: "Medium", value: summary.medium_alerts },
    { name: "Low", value: summary.low_alerts },
  ];

  return (
    <div>
      <h3>Severity Distribution</h3>
      <ResponsiveContainer width="100%" height={250}>
        <BarChart data={data}>
          <XAxis dataKey="name" stroke="#94a3b8" />
          <YAxis stroke="#94a3b8" />
          <Tooltip />
          <Bar dataKey="value" fill="#3b82f6" />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}

export default SeverityChart;