import React, { useEffect, useRef } from "react";
import * as d3 from "d3";

export default function SecurityGraph({ data, onNodeClick }) {
  const svgRef = useRef();

  useEffect(() => {
    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();

    const width = 900;
    const height = 650;

    svg.attr("width", width).attr("height", height);

    const nodes = (data?.nodes || []).map((d) => ({ ...d }));
    const links = (data?.links || []).map((d) => ({ ...d }));

    const simulation = d3
      .forceSimulation(nodes)
      .force(
        "link",
        d3.forceLink(links).id((d) => d.id).distance(140)
      )
      .force("charge", d3.forceManyBody().strength(-500))
      .force("center", d3.forceCenter(width / 2, height / 2));

    const link = svg
      .append("g")
      .selectAll("line")
      .data(links)
      .enter()
      .append("line")
      .attr("stroke", "#94a3b8")
      .attr("stroke-width", 2);

    const linkLabel = svg
      .append("g")
      .selectAll("text")
      .data(links)
      .enter()
      .append("text")
      .text((d) => d.type)
      .attr("font-size", "10px")
      .attr("fill", "#334155");

    const colorByLabel = (label) => {
      switch (label) {
        case "AttackerIP":
          return "#ef4444";
        case "Server":
          return "#3b82f6";
        case "Alert":
          return "#8b5cf6";
        case "APIEndpoint":
          return "#10b981";
        case "ResponseAction":
          return "#6b7280";
        default:
          return "#f59e0b";
      }
    };

    const node = svg
      .append("g")
      .selectAll("circle")
      .data(nodes)
      .enter()
      .append("circle")
      .attr("r", 20)
      .attr("fill", (d) => colorByLabel(d.label))
      .style("cursor", "pointer")
      .call(
        d3
          .drag()
          .on("start", dragStarted)
          .on("drag", dragged)
          .on("end", dragEnded)
      )
      .on("click", function (event, d) {
        if (onNodeClick) {
          onNodeClick(d);
        }
      });

    const label = svg
      .append("g")
      .selectAll("text.node-label")
      .data(nodes)
      .enter()
      .append("text")
      .attr("class", "node-label")
      .text((d) => d.name || d.id)
      .attr("font-size", "12px")
      .attr("text-anchor", "middle")
      .attr("dy", 35);

    simulation.on("tick", () => {
      link
        .attr("x1", (d) => d.source.x)
        .attr("y1", (d) => d.source.y)
        .attr("x2", (d) => d.target.x)
        .attr("y2", (d) => d.target.y);

      linkLabel
        .attr("x", (d) => (d.source.x + d.target.x) / 2)
        .attr("y", (d) => (d.source.y + d.target.y) / 2);

      node
        .attr("cx", (d) => d.x)
        .attr("cy", (d) => d.y);

      label
        .attr("x", (d) => d.x)
        .attr("y", (d) => d.y);
    });

    function dragStarted(event, d) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    }

    function dragged(event, d) {
      d.fx = event.x;
      d.fy = event.y;
    }

    function dragEnded(event, d) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    }

    return () => {
      simulation.stop();
    };
  }, [data, onNodeClick]);

  return (
    <div
      style={{
        border: "1px solid #ccc",
        borderRadius: "8px",
        background: "#f3f4f6",
        overflow: "hidden",
      }}
    >
      <svg ref={svgRef}></svg>
    </div>
  );
}