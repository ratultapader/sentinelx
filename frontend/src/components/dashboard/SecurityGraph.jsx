import React, { useEffect, useState } from "react";
import api from "../../services/api";
import NodeDetailsPanel from "../NodeDetailsPanel";

export default function SecurityGraph({ ip, onNodeClick }) {
  const [graph, setGraph] = useState({
    nodes: [],
    edges: [],
  });

  const [path, setPath] = useState([]);
  const [selected, setSelected] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // ===============================
  // FETCH GRAPH
  // ===============================
  useEffect(() => {
    if (!ip) return;

    setLoading(true);
    setError(null);

    api
      .get(`/api/graph/${encodeURIComponent(ip)}`, {
        headers: { "X-Tenant-ID": "t1" },
      })
      .then((res) => {
        const data = res.data || {};

        setGraph({
          nodes: data.nodes || [],
          edges: data.links || [],
        });
      })
      .catch((err) => {
        console.error("Graph error:", err);

        if (err.response && err.response.status === 404) {
          // NOT an error → just no data
          setGraph({ nodes: [], edges: [] });
        } else {
          setError("Failed to load graph");
        }
      })
      .finally(() => setLoading(false));
  }, [ip]);

  // ===============================
  // FETCH ATTACK PATH
  // ===============================
  useEffect(() => {
    if (!ip) return;

    api
      .get(`/api/attack_path/${encodeURIComponent(ip)}`, {
        headers: { "X-Tenant-ID": "t1" },
      })
      .then((res) => {
        setPath(res.data.path || []);
      })
      .catch((err) => console.error("Path error:", err));
  }, [ip]);

  // ===============================
  // LOADING / ERROR
  // ===============================
  if (loading) {
    return <div style={styles.loading}>Loading graph...</div>;
  }

  if (error) {
    return <div style={styles.error}>{error}</div>;
  }

  return (
    <div style={styles.container}>
      
      {/* ================= GRAPH ================= */}
      <div style={styles.graphArea}>
        <h3>🌐 Attack Graph ({ip})</h3>

        {graph.nodes.length === 0 ? (
          <p style={{ color: "#94a3b8" }}>No graph data found</p>
        ) : (
          graph.nodes.map((node, i) => {
            // ✅ IMPROVED PATH MATCH
            const isInPath =
              path.includes(node.id) ||
              path.includes(node.name) ||
              (node.label &&
                path.some((p) =>
                  p.toLowerCase().includes(node.label.toLowerCase())
                ));

            const isSelected = selected?.id === node.id;

            return (
              <div
                key={node.id || i}
                onClick={() => {
                  console.log("CLICK:", node);

                  setSelected(node);

                  // 🔥 PASS TO PARENT (GraphView)
                  if (onNodeClick) {
                    onNodeClick(node);
                  }
                }}
                onMouseEnter={(e) => {
                  console.log("Hover:", node);
                  e.currentTarget.style.background = "#334155";
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = isSelected
                    ? "#334155"
                    : "#1e293b";
                }}
                style={{
                  ...styles.node,
                  borderLeft: `4px solid ${
                    isInPath
                      ? "#facc15" // 🔥 PATH HIGHLIGHT
                      : getNodeColor(node.label || node.type)
                  }`,
                  background: isSelected ? "#334155" : "#1e293b",
                }}
              >
                <strong>{node.name || node.id}</strong>

                <div style={{ fontSize: "12px", color: "#94a3b8" }}>
                  {node.label || node.type}
                </div>
              </div>
            );
          })
        )}

        {/* ================= EDGES ================= */}
        <h4 style={{ marginTop: "20px" }}>Connections</h4>

        {graph.edges.length === 0 ? (
          <p style={{ color: "#94a3b8" }}>No connections</p>
        ) : (
          graph.edges.map((e, i) => (
            <div key={i} style={styles.edge}>
              {e.source} → {e.target} ({e.type})
            </div>
          ))
        )}
      </div>

      {/* ================= SIDE PANEL ================= */}
      <NodeDetailsPanel node={selected} />
    </div>
  );
}

// ===============================
// STYLES
// ===============================
const styles = {
  container: {
    display: "flex",
    height: "100%",
    color: "white",
  },
  graphArea: {
    flex: 1,
    padding: "12px",
  },
  node: {
    padding: "10px",
    margin: "6px 0",
    cursor: "pointer",
    borderRadius: "6px",
    transition: "0.2s",
  },
  edge: {
    fontSize: "12px",
    color: "#94a3b8",
  },
  loading: {
    padding: "10px",
    color: "white",
  },
  error: {
    padding: "10px",
    color: "#ef4444",
  },
};

// ===============================
// COLOR LOGIC
// ===============================
function getNodeColor(label) {
  if (!label) return "#22c55e";

  const l = label.toLowerCase();

  if (l.includes("alert")) return "#a855f7";
  if (l.includes("response")) return "#94a3b8";
  if (l.includes("api")) return "#3b82f6";
  if (l.includes("attacker") || l.includes("ip")) return "#ef4444";

  return "#22c55e";
}