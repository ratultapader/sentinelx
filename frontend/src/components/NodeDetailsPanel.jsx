import React from "react";

export default function NodeDetailsPanel({ node }) {
  if (!node) return null;

  return (
    <div style={styles.panel}>
      <h3>🔍 Node Details</h3>

      <p><strong>ID:</strong> {node.id}</p>
      <p><strong>Type:</strong> {node.type}</p>
      <p><strong>Label:</strong> {node.label}</p>

      {node.metadata && (
        <pre style={styles.meta}>
          {JSON.stringify(node.metadata, null, 2)}
        </pre>
      )}
    </div>
  );
}

const styles = {
  panel: {
    position: "absolute",
    right: 0,
    top: 0,
    width: "320px",
    height: "100%",
    background: "#0f172a",
    color: "#e2e8f0",
    padding: "16px",
    borderLeft: "1px solid #1e293b",
    overflowY: "auto",
    zIndex: 10
  },
  meta: {
    fontSize: "12px",
    background: "#020617",
    padding: "10px",
    borderRadius: "6px"
  }
};