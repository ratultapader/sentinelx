import React, { useState } from "react";
import { post } from "../services/api";

export default function RuleForm({ onCreated }) {
  const [name, setName] = useState("");
  const [condition, setCondition] = useState("");
  const [action, setAction] = useState("alert");

  const createRule = async () => {
    if (!name || !condition) {
      alert("Fill all fields");
      return;
    }

    await post("/api/rules", { name, condition, action });

    setName("");
    setCondition("");
    setAction("alert");

    onCreated();
  };

  return (
    <div style={{ marginTop: "15px" }}>
      <input
        placeholder="Rule name"
        value={name}
        onChange={(e) => setName(e.target.value)}
      />

      <input
        placeholder="Condition (e.g. req > 50)"
        value={condition}
        onChange={(e) => setCondition(e.target.value)}
      />

      <select value={action} onChange={(e) => setAction(e.target.value)}>
        <option value="alert">Alert</option>
        <option value="block">Block</option>
        <option value="rate_limit">Rate Limit</option>
      </select>

      <button onClick={createRule}>Create</button>
    </div>
  );
}