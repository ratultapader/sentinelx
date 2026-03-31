import React, { useState } from "react";
import SecurityGraph from "./dashboard/SecurityGraph";
import NodeDetailsPanel from "./NodeDetailsPanel";

export default function GraphView({ ip }) {
  const [selectedNode, setSelectedNode] = useState(null);

  return (
    <div style={{ position: "relative", height: "100%" }}>

      {/* GRAPH */}
      <SecurityGraph 
        ip={ip} 
        onNodeClick={(node) => {
          console.log("CLICK NODE:", node);
          setSelectedNode(node);
        }}
      />

      {/* RIGHT PANEL */}
      <NodeDetailsPanel node={selectedNode} />

    </div>
  );
}