export function exportToCSV(data, filename = "data.csv") {
  if (!data || !data.length) return;

  // 🔥 flatten for clean CSV
  const flatData = data.map(item => ({
    id: item.id,
    source_ip: item.source_ip,
    severity: item.severity,
    alert_count: item.alert_count || "",
    timestamp: item.timestamp
  }));

  const headers = Object.keys(flatData[0]).join(",");

  const rows = flatData.map(obj =>
    Object.values(obj)
      .map(val => `"${val}"`)
      .join(",")
  );

  const csv = [headers, ...rows].join("\n");

  const blob = new Blob([csv], { type: "text/csv" });
  const url = URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();

  URL.revokeObjectURL(url);
}

export function exportToJSON(data, filename = "data.json") {
  const blob = new Blob(
    [JSON.stringify(data, null, 2)],
    { type: "application/json" }
  );

  const url = URL.createObjectURL(blob);

  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();

  URL.revokeObjectURL(url);
}