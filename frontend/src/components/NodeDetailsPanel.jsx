import React from "react";

export default function NodeDetailsPanel({ node }) {
  if (!node) {
    return (
      <div
        style={{
          border: "1px solid #ccc",
          borderRadius: "8px",
          padding: "16px",
          background: "#f9f9f9",
          minHeight: "200px",
        }}
      >
        <h3>Node Details</h3>
        <p>Click a node to view details.</p>
      </div>
    );
  }

  return (
    <div
      style={{
        border: "1px solid #ccc",
        borderRadius: "8px",
        padding: "16px",
        background: "#f9f9f9",
        minHeight: "200px",
      }}
    >
      <h3>Node Details</h3>

      <p><strong>ID:</strong> {node.id}</p>
      <p><strong>Label:</strong> {node.label}</p>
      <p><strong>Name:</strong> {node.name}</p>

      <h4>Properties</h4>
      {node.properties && Object.keys(node.properties).length > 0 ? (
        <ul>
          {Object.entries(node.properties).map(([key, value]) => (
            <li key={key}>
              <strong>{key}:</strong> {String(value)}
            </li>
          ))}
        </ul>
      ) : (
        <p>No properties available.</p>
      )}
    </div>
  );
}