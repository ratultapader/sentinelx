import React from "react";

export default function ThreatTable({ ips }) {
  return (
    <table style={styles.table}>
      <thead>
        <tr>
          <th>IP</th>
          <th>Reason</th>
          <th>Source</th>
          <th>Match Count</th>
          <th>Risk</th>
        </tr>
      </thead>

      <tbody>
        {ips.map((i) => (
          <tr key={i.ip} style={highlight(i)}>
            <td>{i.ip}</td>
            <td>{i.reason || "-"}</td>
            <td>{i.source || "-"}</td>
            <td>{i.match_count || 0}</td>
            <td>
              {i.match_count > 10
                ? "🔥 Critical"
                : i.match_count > 5
                ? "⚠ High"
                : "Normal"}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function highlight(i) {
  if (i.match_count > 10) {
    return { background: "#991b1b" };
  }
  if (i.match_count > 5) {
    return { background: "#7f1d1d" };
  }
}

const styles = {
  table: {
    width: "100%",
    background: "#1e293b",
    color: "white",
    borderCollapse: "collapse"
  }
};