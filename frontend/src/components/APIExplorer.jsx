import React, { useState } from "react";
import api from "../services/api";

export default function APIExplorer() {
  const [endpoint, setEndpoint] = useState("/api/alerts");
  const [response, setResponse] = useState("");

  const callAPI = async () => {
    try {
      const res = await api.get(endpoint);
      setResponse(JSON.stringify(res.data, null, 2));
    } catch (err) {
      setResponse("❌ Error calling API");
    }
  };

  return (
    <div style={styles.card}>
      <h3>🔗 API Explorer</h3>

      <input
        value={endpoint}
        onChange={e => setEndpoint(e.target.value)}
        style={styles.input}
      />

      <button onClick={callAPI} style={styles.button}>
        Send
      </button>

      <pre style={styles.output}>{response}</pre>
    </div>
  );
}

const styles = {
  card: {
    background: "#1e293b",
    padding: "15px",
    borderRadius: "10px"
  },
  input: {
    width: "300px",
    marginRight: "10px"
  },
  button: {
    padding: "5px 10px"
  },
  output: {
  marginTop: "10px",
  background: "#020617",
  padding: "10px",
  borderRadius: "5px",
  maxHeight: "250px",
  overflowY: "auto",
  fontSize: "12px"
}
};