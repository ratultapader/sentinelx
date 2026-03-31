import React from "react";

export default function PlaybookTable({ playbooks, onRefresh }) {

  const toggle = (p) => {
    fetch(`http://localhost:9090/api/playbooks/${p.id}/toggle`, {
      method: "POST",
      headers: { "X-Tenant-ID": "t1" }
    }).then(onRefresh);
  };

  const remove = (id) => {
    fetch(`http://localhost:9090/api/playbooks/${id}`, {
      method: "DELETE",
      headers: { "X-Tenant-ID": "t1" }
    }).then(onRefresh);
  };

  return (
    <table style={styles.table}>
      <thead>
        <tr>
          <th>#</th>
          <th>Condition</th>
          <th>Action</th>
          <th>Enabled</th>
          <th>Delete</th>
        </tr>
      </thead>

      <tbody>
        {playbooks.map((p, index) => (
          <tr key={p.id} style={styles.row}>
            <td>{index + 1}</td>

            <td>
              <strong>IF</strong> {p.condition}
            </td>

            <td>
              <strong>THEN</strong> {p.action}
            </td>

            <td>
              <input
                type="checkbox"
                checked={p.enabled}
                onChange={() => toggle(p)}
              />
            </td>

            <td>
              <button onClick={() => remove(p.id)} style={styles.deleteBtn}>
                ❌
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

const styles = {
  table: {
    width: "100%",
    background: "#1e293b",
    color: "white",
    borderCollapse: "collapse"
  },
  row: {
    borderBottom: "1px solid #334155"
  },
  deleteBtn: {
    background: "red",
    color: "white",
    border: "none",
    padding: "5px",
    cursor: "pointer"
  }
};