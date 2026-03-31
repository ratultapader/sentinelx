import React, { useEffect, useRef, useState } from "react";
import LiveFeed from "../components/dashboard/LiveFeed";

export default function Live() {
  const [events, setEvents] = useState([]);
  const [paused, setPaused] = useState(false);

  const wsRef = useRef(null);
  const reconnectRef = useRef(null);
  const seenEvents = useRef(new Set()); // 🔥 deduplication

  useEffect(() => {
    // ✅ prevent multiple connections (React Strict Mode fix)
    if (wsRef.current) return;

    connectWS();

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }

      if (reconnectRef.current) {
        clearTimeout(reconnectRef.current);
      }
    };
  }, []);

  const connectWS = () => {
    // ✅ prevent duplicate connection
    if (wsRef.current) {
      console.log("⚠️ WS already exists, skipping...");
      return;
    }

    console.log("🔌 Connecting WebSocket...");

    const ws = new WebSocket("ws://localhost:9090/ws?tenant_id=t1");
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("✅ WebSocket connected");
    };

    ws.onmessage = (msg) => {
      try {
        const data = JSON.parse(msg.data);

        // ✅ ignore system messages
        if (data.type === "connected" || data.type === "ping") return;

        // 🔥 create unique key (VERY IMPORTANT)
        const uniqueKey =
          data.id ||
          `${data.source_ip}-${data.timestamp}-${data.type}`;

        // ❌ skip duplicates
        if (seenEvents.current.has(uniqueKey)) return;

        seenEvents.current.add(uniqueKey);

        setEvents((prev) => {
          if (paused) return prev;
          return [data, ...prev.slice(0, 50)];
        });

      } catch (e) {
        console.error("Invalid WS data", e);
      }
    };

    ws.onclose = () => {
      console.warn("⚠️ WS disconnected");

      wsRef.current = null;

      // ✅ prevent multiple reconnect timers
      if (!reconnectRef.current) {
        reconnectRef.current = setTimeout(() => {
          reconnectRef.current = null;
          connectWS();
        }, 2000);
      }
    };

    ws.onerror = () => {
      // optional: silence browser warning
    };
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1 style={{ color: "white" }}>🔴 Live Monitoring</h1>

      <div style={{ marginBottom: "10px" }}>
        <button onClick={() => setPaused((p) => !p)}>
          {paused ? "Resume" : "Pause"}
        </button>

        <button
          onClick={() => {
            setEvents([]);
            seenEvents.current.clear(); // 🔥 reset dedup cache
          }}
          style={{ marginLeft: "10px" }}
        >
          Clear
        </button>
      </div>

      <LiveFeed events={events} />
    </div>
  );
}