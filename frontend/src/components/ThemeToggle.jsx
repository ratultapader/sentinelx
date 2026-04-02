import React, { useEffect, useState } from "react";

export default function ThemeToggle() {
  const [theme, setTheme] = useState("dark");

  useEffect(() => {
    const saved = localStorage.getItem("theme");

    if (saved) {
      setTheme(saved);
      document.body.className = saved;
    } else {
      document.body.className = "dark";
    }
  }, []);

  const toggle = () => {
    const newTheme = theme === "dark" ? "light" : "dark";

    setTheme(newTheme);
    localStorage.setItem("theme", newTheme);
    document.body.className = newTheme;
  };

  return (
    <button onClick={toggle}>
      {theme === "dark" ? "☀ Light" : "🌙 Dark"}
    </button>
  );
}