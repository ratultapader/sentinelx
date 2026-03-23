import React, { useState } from "react";
import SecurityGraph from "../components/SecurityGraph";
import NodeDetailsPanel from "../components/NodeDetailsPanel";
import { fetchGraphBySourceIP } from "../services/graphApi";

export default function GraphInvestigationPage() {
  const [sourceIP, setSourceIP] = useState("192.168.1.5");
  const [graphData, setGraphData] = useState({ nodes: [], links: [] });
  const [selectedNode, setSelectedNode] = useState(null);
  const [error, setError] = useState("");

  const loadGraph = async () => {
    try {
      setError("");
      const data = await fetchGraphBySourceIP(sourceIP);
      setGraphData(data);
      setSelectedNode(null);
    } catch (err) {
      setError(err.message || "Failed to fetch");
    }
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1>Security Graph Investigation</h1>

      <div style={{ marginBottom: "16px" }}>
        <input
          value={sourceIP}
          onChange={(e) => setSourceIP(e.target.value)}
          placeholder="Enter source IP"
        />
        <button onClick={loadGraph} style={{ marginLeft: "8px" }}>
          Load Graph
        </button>
      </div>

      {error && <p style={{ color: "red" }}>{error}</p>}

      <div
        style={{
          display: "flex",
          gap: "20px",
          alignItems: "flex-start",
        }}
      >
        <div style={{ flex: 2 }}>
          <SecurityGraph data={graphData} onNodeClick={setSelectedNode} />
        </div>

        <div style={{ flex: 1 }}>
          <NodeDetailsPanel node={selectedNode} />
        </div>
      </div>
    </div>
  );
}