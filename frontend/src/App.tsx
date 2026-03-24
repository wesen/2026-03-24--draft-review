import { MacWindow, MacButton, MenuBar } from "./chrome";
import { ProgressBar, StatCard } from "./primitives";

const menus = [
  {
    label: "File",
    items: [
      { label: "New Article\u2026", shortcut: "\u2318N" },
      { divider: true, label: "" },
      { label: "Export Feedback\u2026", shortcut: "\u2318E" },
    ],
  },
  { label: "Edit", items: [{ label: "Undo", shortcut: "\u2318Z" }] },
  { label: "View", items: [{ label: "Dashboard" }] },
  { label: "Help", items: [{ label: "Beta Reading Guide" }] },
];

export default function App() {
  return (
    <div className="dr-desktop">
      <MenuBar menus={menus} rightStatus="0 reactions" />
      <MacWindow title="Draft Review \u2014 Dashboard" maximized>
        <div style={{ padding: 16 }}>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "repeat(4, 1fr)",
              gap: 10,
              marginBottom: 16,
            }}
          >
            <StatCard icon="\u25C9" value={5} label="Readers" />
            <StatCard icon="\u2726" value={10} label="Reactions" />
            <StatCard icon="\u00A7" value={5} label="Sections" />
            <StatCard icon="\u25B8" value="60%" label="Avg Progress" />
          </div>
          <div style={{ maxWidth: 300, marginBottom: 16 }}>
            <ProgressBar percent={60} />
          </div>
          <MacButton primary>Open Full Review</MacButton>
        </div>
      </MacWindow>
    </div>
  );
}
