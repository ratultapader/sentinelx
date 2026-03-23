export async function fetchGraphBySourceIP(sourceIP) {
  const res = await fetch(
    `http://localhost:9090/api/graph?source_ip=${encodeURIComponent(sourceIP)}`
  );

  if (!res.ok) {
    throw new Error(`failed to fetch graph: ${res.status}`);
  }

  return res.json();
}