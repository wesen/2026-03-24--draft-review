import { useState, useEffect, useRef, useCallback } from "react";

const CHICAGO = `"Chicago_12", "ChicagoFLF", "Geneva", monospace`;
const GENEVA = `"Geneva", "Monaco", monospace`;

// ─── pixel‑perfect palette ───
const P = {
  bg: "#e8e8e8",
  white: "#ffffff",
  black: "#000000",
  gray: "#a0a0a0",
  darkGray: "#555555",
  lightGray: "#c0c0c0",
  stripes: "repeating-linear-gradient(0deg,#fff 0px,#fff 1px,#000 1px,#000 2px)",
};

// ─── sample data ───
const ARTICLES = [
  {
    id: 1,
    title: "Why Design Systems Fail",
    status: "In Review",
    sections: [
      { id: "s1", title: "Introduction", text: "Design systems promise consistency and speed, but many teams abandon them within a year. This article explores why that happens and what you can do about it. The problem isn't usually technical — it's organizational." },
      { id: "s2", title: "The Adoption Cliff", text: "Most design systems see strong initial adoption. Engineers are excited, designers feel heard, and leadership is optimistic. But around month 6, usage plateaus. Teams start creating workarounds, one-off components multiply, and the system starts to feel like a constraint rather than an enabler." },
      { id: "s3", title: "Governance Gaps", text: "Without clear ownership, design systems drift. Who decides when to add a new component? Who approves breaking changes? When these questions go unanswered, teams lose trust in the system and start going rogue." },
      { id: "s4", title: "The Documentation Problem", text: "Even well-built systems fail without good docs. If a developer can't figure out how to use a component in under 2 minutes, they'll build their own. Documentation isn't a nice-to-have — it's the product." },
      { id: "s5", title: "Making It Work", text: "The teams that succeed treat their design system like a product, not a project. They have dedicated maintainers, regular release cycles, feedback loops with consumers, and executive sponsorship that actually means something." },
    ],
    readers: [
      { id: "r1", name: "Sarah K.", progress: 100, avatar: "SK" },
      { id: "r2", name: "Marcus T.", progress: 60, avatar: "MT" },
      { id: "r3", name: "Priya R.", progress: 80, avatar: "PR" },
      { id: "r4", name: "James L.", progress: 40, avatar: "JL" },
      { id: "r5", name: "Chen W.", progress: 20, avatar: "CW" },
    ],
    reactions: [
      { section: "s1", reader: "Sarah K.", type: "useful", text: "Great hook — immediately relevant to my situation." },
      { section: "s2", reader: "Marcus T.", type: "confused", text: "What do you mean by 'workarounds'? Can you give an example?" },
      { section: "s2", reader: "Priya R.", type: "slow", text: "This section drags a bit. Maybe tighten the middle paragraph?" },
      { section: "s2", reader: "Sarah K.", type: "favorite", text: "YES. This is exactly what happened at my last company." },
      { section: "s3", reader: "James L.", type: "confused", text: "Who typically owns this? Would help to see a real org chart." },
      { section: "s3", reader: "Priya R.", type: "useful", text: "The governance framework idea is gold. Worth expanding." },
      { section: "s4", reader: "Sarah K.", type: "useful", text: "The 2-minute rule is a great benchmark." },
      { section: "s4", reader: "Marcus T.", type: "favorite", text: "Documentation IS the product. Quotable." },
      { section: "s5", reader: "Priya R.", type: "slow", text: "Feels rushed compared to earlier sections. Needs more depth." },
      { section: "s1", reader: "Chen W.", type: "useful", text: "Strong opening. Sets the right expectations." },
    ],
  },
  {
    id: 2,
    title: "Remote Work Isn't Working",
    status: "Draft",
    sections: [
      { id: "s1", title: "The Broken Promise", text: "We were told remote work would free us. In many ways it has. But the cracks are showing: loneliness, miscommunication, and a creeping sense that we've traded one set of problems for another." },
      { id: "s2", title: "What We Lost", text: "Serendipity. The hallway conversation that sparks an idea. The lunch that becomes a brainstorm. These aren't just perks — they're how organizations actually think." },
      { id: "s3", title: "A Middle Path", text: "The answer isn't going back to the office five days a week. It's being intentional about when and why we gather. Async by default, sync with purpose." },
    ],
    readers: [
      { id: "r1", name: "Sarah K.", progress: 100, avatar: "SK" },
      { id: "r2", name: "James L.", progress: 33, avatar: "JL" },
    ],
    reactions: [
      { section: "s1", reader: "Sarah K.", type: "favorite", text: "This opening paragraph is perfect." },
      { section: "s1", reader: "James L.", type: "confused", text: "What cracks specifically? Feels vague." },
    ],
  },
  {
    id: 3,
    title: "The Art of Saying No",
    status: "New",
    sections: [
      { id: "s1", title: "Why We Say Yes", text: "Saying yes feels safe. It keeps the peace, maintains relationships, and avoids the discomfort of disappointing someone. But every yes to something unimportant is a no to something that matters." },
    ],
    readers: [],
    reactions: [],
  },
];

const REACTION_ICONS = {
  useful: { icon: "★", label: "Useful", bg: "#000" },
  confused: { icon: "?", label: "Confused", bg: "#000" },
  slow: { icon: "◎", label: "Slow", bg: "#000" },
  favorite: { icon: "♥", label: "Favorite", bg: "#000" },
};

// ─── Mac OS 1 Window chrome ───
function MacWindow({ title, children, x = 40, y = 40, w = 600, h = 450, onClose, zIndex = 1, style = {}, maximized = false }) {
  const [pos, setPos] = useState({ x, y });
  const [dragging, setDragging] = useState(false);
  const dragRef = useRef(null);

  const onMouseDown = (e) => {
    if (maximized) return;
    dragRef.current = { startX: e.clientX - pos.x, startY: e.clientY - pos.y };
    setDragging(true);
  };
  useEffect(() => {
    if (!dragging) return;
    const onMove = (e) => {
      if (dragRef.current) setPos({ x: e.clientX - dragRef.current.startX, y: e.clientY - dragRef.current.startY });
    };
    const onUp = () => setDragging(false);
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => { window.removeEventListener("mousemove", onMove); window.removeEventListener("mouseup", onUp); };
  }, [dragging]);

  const winStyle = maximized
    ? { position: "absolute", left: 0, top: 22, right: 0, bottom: 0, zIndex, display: "flex", flexDirection: "column" }
    : { position: "absolute", left: pos.x, top: pos.y, width: w, height: h, zIndex, display: "flex", flexDirection: "column", ...style };

  return (
    <div style={{
      ...winStyle,
      border: `2px solid ${P.black}`,
      background: P.white,
      boxShadow: "3px 3px 0 rgba(0,0,0,0.3)",
      borderRadius: 3,
      overflow: "hidden",
    }}>
      {/* Title bar */}
      <div onMouseDown={onMouseDown} style={{
        background: P.white,
        borderBottom: `2px solid ${P.black}`,
        padding: "2px 6px",
        display: "flex",
        alignItems: "center",
        cursor: maximized ? "default" : "grab",
        minHeight: 22,
        userSelect: "none",
        flexShrink: 0,
      }}>
        {/* Close box */}
        <div onClick={onClose} style={{
          width: 14, height: 14,
          border: `1.5px solid ${P.black}`,
          borderRadius: 1,
          cursor: "pointer",
          flexShrink: 0,
          display: "flex", alignItems: "center", justifyContent: "center",
          marginRight: 8,
        }}>
          <span style={{ fontSize: 9, lineHeight: 1 }}></span>
        </div>
        {/* Stripes left */}
        <div style={{ flex: 1, height: 10, background: P.stripes, marginRight: 8 }} />
        <span style={{ fontFamily: CHICAGO, fontSize: 12, fontWeight: "bold", whiteSpace: "nowrap", letterSpacing: 0.3 }}>{title}</span>
        <div style={{ flex: 1, height: 10, background: P.stripes, marginLeft: 8 }} />
      </div>
      {/* Content */}
      <div style={{ flex: 1, overflow: "auto", position: "relative" }}>
        {children}
      </div>
    </div>
  );
}

// ─── Retro Button ───
function MacButton({ children, onClick, primary = false, style = {}, disabled = false }) {
  return (
    <button disabled={disabled} onClick={onClick} style={{
      fontFamily: CHICAGO, fontSize: 11, padding: "4px 16px",
      border: `2px solid ${P.black}`,
      borderRadius: 8,
      background: primary ? P.black : P.white,
      color: primary ? P.white : P.black,
      cursor: disabled ? "default" : "pointer",
      boxShadow: primary ? "inset 0 0 0 1px #fff, 0 0 0 2px #000" : "none",
      opacity: disabled ? 0.4 : 1,
      ...style,
    }}>
      {children}
    </button>
  );
}

// ─── Menu Bar ───
function MenuBar({ onNewArticle, onAbout, currentView, setView }) {
  const [openMenu, setOpenMenu] = useState(null);

  return (
    <div style={{
      position: "fixed", top: 0, left: 0, right: 0, height: 22,
      background: P.white, borderBottom: `2px solid ${P.black}`,
      display: "flex", alignItems: "center", zIndex: 9999,
      fontFamily: CHICAGO, fontSize: 12,
    }}>
      {/* Apple logo */}
      <div style={{ position: "relative" }}
        onMouseEnter={() => openMenu && setOpenMenu("apple")}
        onClick={() => setOpenMenu(openMenu === "apple" ? null : "apple")}>
        <div style={{ padding: "0 12px", cursor: "pointer", fontWeight: "bold", fontSize: 14 }}></div>
        {openMenu === "apple" && (
          <MenuDropdown items={[
            { label: "About Draft Review", action: onAbout },
            { divider: true },
            { label: "Version 1.0" },
          ]} onClose={() => setOpenMenu(null)} />
        )}
      </div>
      <div style={{ position: "relative" }}
        onMouseEnter={() => openMenu && setOpenMenu("file")}
        onClick={() => setOpenMenu(openMenu === "file" ? null : "file")}>
        <div style={{ padding: "0 10px", cursor: "pointer", fontWeight: openMenu === "file" ? "bold" : "normal" }}>File</div>
        {openMenu === "file" && (
          <MenuDropdown items={[
            { label: "New Article…", action: onNewArticle, shortcut: "⌘N" },
            { label: "Import from Docs…", shortcut: "⌘I" },
            { divider: true },
            { label: "Export Feedback…", shortcut: "⌘E" },
          ]} onClose={() => setOpenMenu(null)} />
        )}
      </div>
      <div style={{ position: "relative" }}
        onMouseEnter={() => openMenu && setOpenMenu("view")}
        onClick={() => setOpenMenu(openMenu === "view" ? null : "view")}>
        <div style={{ padding: "0 10px", cursor: "pointer" }}>View</div>
        {openMenu === "view" && (
          <MenuDropdown items={[
            { label: "✓ Dashboard", action: () => { setView("dashboard"); setOpenMenu(null); } },
            { label: "  Reader List", action: () => setOpenMenu(null) },
            { divider: true },
            { label: "  Show Reactions", action: () => setOpenMenu(null) },
          ]} onClose={() => setOpenMenu(null)} />
        )}
      </div>
      <div style={{ position: "relative" }}
        onMouseEnter={() => openMenu && setOpenMenu("help")}
        onClick={() => setOpenMenu(openMenu === "help" ? null : "help")}>
        <div style={{ padding: "0 10px", cursor: "pointer" }}>Help</div>
        {openMenu === "help" && (
          <MenuDropdown items={[
            { label: "Beta Reading Guide" },
            { label: "How to Invite Readers" },
            { divider: true },
            { label: "Keyboard Shortcuts" },
          ]} onClose={() => setOpenMenu(null)} />
        )}
      </div>
      {openMenu && <div style={{ position: "fixed", inset: 0, zIndex: -1 }} onClick={() => setOpenMenu(null)} />}
    </div>
  );
}

function MenuDropdown({ items, onClose }) {
  return (
    <div style={{
      position: "absolute", top: 22, left: 0, minWidth: 200,
      background: P.white, border: `2px solid ${P.black}`,
      boxShadow: "3px 3px 0 rgba(0,0,0,0.25)",
      zIndex: 10000, padding: "2px 0",
    }}>
      {items.map((item, i) =>
        item.divider ? (
          <div key={i} style={{ borderTop: `1px solid ${P.black}`, margin: "2px 0" }} />
        ) : (
          <div key={i} onClick={() => { item.action?.(); onClose(); }} style={{
            padding: "3px 20px", cursor: item.action ? "pointer" : "default",
            fontFamily: CHICAGO, fontSize: 12, display: "flex", justifyContent: "space-between",
            whiteSpace: "nowrap",
          }}
            onMouseEnter={(e) => { e.currentTarget.style.background = P.black; e.currentTarget.style.color = P.white; }}
            onMouseLeave={(e) => { e.currentTarget.style.background = P.white; e.currentTarget.style.color = P.black; }}
          >
            <span>{item.label}</span>
            {item.shortcut && <span style={{ marginLeft: 30, opacity: 0.6 }}>{item.shortcut}</span>}
          </div>
        )
      )}
    </div>
  );
}

// ─── Progress Bar (retro style) ───
function ProgressBar({ percent, width = 120, height = 12 }) {
  return (
    <div style={{
      width, height, border: `1.5px solid ${P.black}`, background: P.white,
      position: "relative", borderRadius: 1,
    }}>
      <div style={{
        position: "absolute", left: 0, top: 0, bottom: 0,
        width: `${percent}%`,
        background: `repeating-linear-gradient(45deg, ${P.black} 0px, ${P.black} 2px, ${P.white} 2px, ${P.white} 4px)`,
      }} />
    </div>
  );
}

// ─── Dashboard view ───
function Dashboard({ articles, onSelectArticle, onInvite }) {
  const [selectedId, setSelectedId] = useState(articles[0]?.id);
  const selected = articles.find((a) => a.id === selectedId) || articles[0];

  const totalReactions = selected.reactions.length;
  const reactionCounts = {};
  Object.keys(REACTION_ICONS).forEach((k) => { reactionCounts[k] = selected.reactions.filter((r) => r.type === k).length; });
  const sectionReactions = {};
  selected.sections.forEach((s) => {
    sectionReactions[s.id] = selected.reactions.filter((r) => r.section === s.id);
  });

  // Find the "book killer" — section with most confused + slow
  let bookKiller = null;
  let maxPain = 0;
  selected.sections.forEach((s) => {
    const pain = (sectionReactions[s.id] || []).filter((r) => r.type === "confused" || r.type === "slow").length;
    if (pain > maxPain) { maxPain = pain; bookKiller = s; }
  });

  return (
    <div style={{ padding: 16, fontFamily: GENEVA, fontSize: 12 }}>
      {/* Article selector */}
      <div style={{ display: "flex", gap: 8, marginBottom: 16, flexWrap: "wrap" }}>
        {articles.map((a) => (
          <div key={a.id} onClick={() => setSelectedId(a.id)} style={{
            padding: "6px 14px", border: `2px solid ${P.black}`,
            background: a.id === selectedId ? P.black : P.white,
            color: a.id === selectedId ? P.white : P.black,
            cursor: "pointer", fontFamily: CHICAGO, fontSize: 11, borderRadius: 2,
            display: "flex", gap: 8, alignItems: "center",
          }}>
            <span style={{ fontSize: 8, padding: "1px 5px", border: `1px solid ${a.id === selectedId ? P.white : P.black}`, borderRadius: 2 }}>
              {a.status}
            </span>
            {a.title}
          </div>
        ))}
      </div>

      {/* Stats row */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(130px, 1fr))", gap: 10, marginBottom: 16 }}>
        {[
          { label: "Readers", value: selected.readers.length, icon: "◉" },
          { label: "Reactions", value: totalReactions, icon: "✦" },
          { label: "Sections", value: selected.sections.length, icon: "§" },
          { label: "Avg Progress", value: selected.readers.length ? Math.round(selected.readers.reduce((a, r) => a + r.progress, 0) / selected.readers.length) + "%" : "—", icon: "▸" },
        ].map((stat) => (
          <div key={stat.label} style={{
            border: `2px solid ${P.black}`, padding: "10px 12px", background: P.white, borderRadius: 2,
          }}>
            <div style={{ fontFamily: CHICAGO, fontSize: 20, marginBottom: 2 }}>{stat.icon} {stat.value}</div>
            <div style={{ fontSize: 10, color: P.darkGray, textTransform: "uppercase", letterSpacing: 1 }}>{stat.label}</div>
          </div>
        ))}
      </div>

      {/* Main grid: readers + reactions graph */}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12, marginBottom: 16 }}>
        {/* Readers panel */}
        <div style={{ border: `2px solid ${P.black}`, borderRadius: 2 }}>
          <div style={{
            padding: "4px 10px", borderBottom: `2px solid ${P.black}`,
            fontFamily: CHICAGO, fontSize: 11, display: "flex", justifyContent: "space-between", alignItems: "center",
          }}>
            <span>Readers</span>
            <MacButton onClick={onInvite} style={{ fontSize: 10, padding: "1px 10px" }}>+ Invite</MacButton>
          </div>
          <div style={{ padding: 8 }}>
            {selected.readers.length === 0 ? (
              <div style={{ padding: 20, textAlign: "center", color: P.gray }}>
                <div style={{ fontSize: 28, marginBottom: 6 }}>📭</div>
                No readers yet. Invite your first beta reader!
              </div>
            ) : selected.readers.map((r) => (
              <div key={r.id} style={{
                display: "flex", alignItems: "center", gap: 8, padding: "5px 4px",
                borderBottom: `1px dotted ${P.lightGray}`,
              }}>
                <div style={{
                  width: 24, height: 24, borderRadius: "50%", border: `1.5px solid ${P.black}`,
                  display: "flex", alignItems: "center", justifyContent: "center",
                  fontFamily: CHICAGO, fontSize: 8, fontWeight: "bold", flexShrink: 0,
                }}>{r.avatar}</div>
                <span style={{ flex: 1, fontSize: 11 }}>{r.name}</span>
                <ProgressBar percent={r.progress} width={80} height={10} />
                <span style={{ fontSize: 10, width: 32, textAlign: "right", color: P.darkGray }}>{r.progress}%</span>
              </div>
            ))}
          </div>
        </div>

        {/* Reactions per section (bar chart) */}
        <div style={{ border: `2px solid ${P.black}`, borderRadius: 2 }}>
          <div style={{ padding: "4px 10px", borderBottom: `2px solid ${P.black}`, fontFamily: CHICAGO, fontSize: 11 }}>
            Reactions by Section
          </div>
          <div style={{ padding: 10 }}>
            {selected.sections.map((s) => {
              const sr = sectionReactions[s.id] || [];
              const isKiller = bookKiller && s.id === bookKiller.id;
              return (
                <div key={s.id} style={{ marginBottom: 8 }}>
                  <div style={{ display: "flex", justifyContent: "space-between", fontSize: 10, marginBottom: 3 }}>
                    <span style={{ fontWeight: isKiller ? "bold" : "normal" }}>
                      {isKiller && "⚠ "}{s.title}
                    </span>
                    <span style={{ color: P.darkGray }}>{sr.length} reactions</span>
                  </div>
                  <div style={{ display: "flex", gap: 2, height: 14 }}>
                    {Object.entries(REACTION_ICONS).map(([type, info]) => {
                      const c = sr.filter((r) => r.type === type).length;
                      if (c === 0) return null;
                      const patterns = {
                        useful: P.black,
                        confused: `repeating-linear-gradient(45deg,${P.black} 0px,${P.black} 2px,${P.white} 2px,${P.white} 4px)`,
                        slow: `repeating-linear-gradient(0deg,${P.black} 0px,${P.black} 1px,${P.white} 1px,${P.white} 3px)`,
                        favorite: `repeating-linear-gradient(90deg,${P.black} 0px,${P.black} 3px,${P.white} 3px,${P.white} 5px)`,
                      };
                      return (
                        <div key={type} title={`${c} ${info.label}`} style={{
                          flex: c, height: "100%", background: patterns[type],
                          border: `1px solid ${P.black}`,
                        }} />
                      );
                    })}
                    {sr.length === 0 && <div style={{ flex: 1, height: "100%", border: `1px dashed ${P.lightGray}` }} />}
                  </div>
                </div>
              );
            })}
            {/* Legend */}
            <div style={{ display: "flex", gap: 10, marginTop: 10, fontSize: 9, color: P.darkGray, flexWrap: "wrap" }}>
              {Object.entries(REACTION_ICONS).map(([type, info]) => (
                <span key={type}>{info.icon} {info.label}</span>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Book Killer alert */}
      {bookKiller && maxPain > 0 && (
        <div style={{
          border: `2px solid ${P.black}`, padding: 12, marginBottom: 16, borderRadius: 2,
          background: `repeating-linear-gradient(45deg, ${P.white} 0px, ${P.white} 8px, #f5f5f5 8px, #f5f5f5 16px)`,
        }}>
          <div style={{ fontFamily: CHICAGO, fontSize: 12, marginBottom: 4 }}>
            ⚠ Potential Draft-Killer Detected
          </div>
          <div style={{ fontSize: 11, lineHeight: 1.5 }}>
            <strong>"{bookKiller.title}"</strong> has the most confusion and pacing issues ({maxPain} flags).
            Readers are telling you this section needs work. Focus your next revision here.
          </div>
          <MacButton onClick={() => onSelectArticle(selected.id, bookKiller.id)} style={{ marginTop: 8, fontSize: 10 }}>
            Review This Section →
          </MacButton>
        </div>
      )}

      {/* Recent feedback */}
      <div style={{ border: `2px solid ${P.black}`, borderRadius: 2 }}>
        <div style={{ padding: "4px 10px", borderBottom: `2px solid ${P.black}`, fontFamily: CHICAGO, fontSize: 11 }}>
          Recent Feedback
        </div>
        <div style={{ maxHeight: 200, overflow: "auto" }}>
          {selected.reactions.slice(-6).reverse().map((r, i) => (
            <div key={i} style={{
              padding: "8px 12px", borderBottom: `1px dotted ${P.lightGray}`,
              display: "flex", gap: 10, alignItems: "flex-start",
            }}>
              <span style={{ fontSize: 16, flexShrink: 0 }}>{REACTION_ICONS[r.type].icon}</span>
              <div>
                <div style={{ fontSize: 10, color: P.darkGray, marginBottom: 2 }}>
                  <strong>{r.reader}</strong> on <em>{selected.sections.find((s) => s.id === r.section)?.title}</em>
                </div>
                <div style={{ fontSize: 11 }}>{r.text}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      <div style={{ textAlign: "center", marginTop: 16 }}>
        <MacButton primary onClick={() => onSelectArticle(selected.id)}>
          Open Full Review →
        </MacButton>
      </div>
    </div>
  );
}

// ─── Article reader view with inline reactions ───
function ArticleReader({ article, onBack, focusSection }) {
  const [selectedSection, setSelectedSection] = useState(focusSection || article.sections[0]?.id);
  const [showInvite, setShowInvite] = useState(false);
  const [reactionFilter, setReactionFilter] = useState(null);
  const section = article.sections.find((s) => s.id === selectedSection);
  const sectionReactions = article.reactions.filter((r) => r.section === selectedSection);
  const filtered = reactionFilter ? sectionReactions.filter((r) => r.type === reactionFilter) : sectionReactions;

  return (
    <div style={{ display: "flex", height: "100%", fontFamily: GENEVA, fontSize: 12 }}>
      {/* Sidebar — section list */}
      <div style={{
        width: 180, borderRight: `2px solid ${P.black}`, display: "flex", flexDirection: "column",
        background: P.white, flexShrink: 0,
      }}>
        <div style={{ padding: "8px 10px", borderBottom: `1px solid ${P.black}`, fontFamily: CHICAGO, fontSize: 10, textTransform: "uppercase", letterSpacing: 1, color: P.darkGray }}>
          Sections
        </div>
        <div style={{ flex: 1, overflow: "auto" }}>
          {article.sections.map((s, i) => {
            const sr = article.reactions.filter((r) => r.section === s.id);
            const isActive = s.id === selectedSection;
            const painCount = sr.filter((r) => r.type === "confused" || r.type === "slow").length;
            return (
              <div key={s.id} onClick={() => { setSelectedSection(s.id); setReactionFilter(null); }} style={{
                padding: "8px 10px", cursor: "pointer",
                background: isActive ? P.black : "transparent",
                color: isActive ? P.white : P.black,
                borderBottom: `1px dotted ${P.lightGray}`,
              }}>
                <div style={{ fontSize: 11, fontWeight: "bold", marginBottom: 3 }}>{i + 1}. {s.title}</div>
                <div style={{ display: "flex", gap: 6, fontSize: 10 }}>
                  {Object.entries(REACTION_ICONS).map(([type, info]) => {
                    const c = sr.filter((r) => r.type === type).length;
                    return c > 0 ? <span key={type} style={{ opacity: isActive ? 0.8 : 0.6 }}>{info.icon}{c}</span> : null;
                  })}
                  {painCount >= 2 && <span>⚠</span>}
                </div>
              </div>
            );
          })}
        </div>
        <div style={{ padding: 8, borderTop: `1px solid ${P.black}` }}>
          <MacButton onClick={onBack} style={{ fontSize: 10, width: "100%" }}>
            ← Dashboard
          </MacButton>
        </div>
      </div>

      {/* Main content */}
      <div style={{ flex: 1, display: "flex", flexDirection: "column" }}>
        {/* Article text */}
        <div style={{
          flex: 1, padding: "20px 28px", overflow: "auto",
          borderBottom: `2px solid ${P.black}`,
        }}>
          <h2 style={{ fontFamily: CHICAGO, fontSize: 16, margin: "0 0 4px", letterSpacing: -0.3 }}>
            {section?.title}
          </h2>
          <div style={{ fontSize: 9, color: P.darkGray, marginBottom: 16, textTransform: "uppercase", letterSpacing: 1 }}>
            Section {article.sections.findIndex((s) => s.id === selectedSection) + 1} of {article.sections.length}
          </div>
          <p style={{ lineHeight: 1.8, fontSize: 13, maxWidth: 540, margin: 0 }}>
            {section?.text}
          </p>

          {/* Reader progress for this section */}
          <div style={{ marginTop: 24, padding: "10px 12px", border: `1px solid ${P.lightGray}`, borderRadius: 2 }}>
            <div style={{ fontSize: 10, color: P.darkGray, marginBottom: 6, fontFamily: CHICAGO }}>Reader Progress Through Here</div>
            <div style={{ display: "flex", gap: 6, flexWrap: "wrap" }}>
              {article.readers.map((r) => {
                const sectionIndex = article.sections.findIndex((s) => s.id === selectedSection);
                const readerSectionProgress = Math.floor((r.progress / 100) * article.sections.length);
                const reachedHere = readerSectionProgress > sectionIndex;
                return (
                  <div key={r.id} style={{
                    display: "flex", alignItems: "center", gap: 4, padding: "2px 6px",
                    border: `1px solid ${P.black}`, borderRadius: 8, fontSize: 10,
                    opacity: reachedHere ? 1 : 0.35,
                  }}>
                    <span style={{ fontSize: 8 }}>{reachedHere ? "●" : "○"}</span>
                    {r.name}
                  </div>
                );
              })}
              {article.readers.length === 0 && <span style={{ color: P.gray, fontSize: 10 }}>No readers yet</span>}
            </div>
          </div>
        </div>

        {/* Reactions panel */}
        <div style={{ height: 220, display: "flex", flexDirection: "column", flexShrink: 0 }}>
          {/* Reaction filter tabs */}
          <div style={{
            display: "flex", borderBottom: `1px solid ${P.black}`, fontFamily: CHICAGO, fontSize: 10, flexShrink: 0,
          }}>
            <div onClick={() => setReactionFilter(null)} style={{
              padding: "5px 12px", cursor: "pointer",
              borderRight: `1px solid ${P.black}`,
              background: !reactionFilter ? P.black : P.white,
              color: !reactionFilter ? P.white : P.black,
            }}>All ({sectionReactions.length})</div>
            {Object.entries(REACTION_ICONS).map(([type, info]) => {
              const c = sectionReactions.filter((r) => r.type === type).length;
              return (
                <div key={type} onClick={() => setReactionFilter(reactionFilter === type ? null : type)} style={{
                  padding: "5px 10px", cursor: "pointer",
                  borderRight: `1px solid ${P.black}`,
                  background: reactionFilter === type ? P.black : P.white,
                  color: reactionFilter === type ? P.white : P.black,
                }}>
                  {info.icon} {info.label} ({c})
                </div>
              );
            })}
          </div>
          {/* Reaction list */}
          <div style={{ flex: 1, overflow: "auto", padding: 8 }}>
            {filtered.length === 0 ? (
              <div style={{ padding: 16, textAlign: "center", color: P.gray }}>
                No {reactionFilter || ""} reactions for this section yet.
              </div>
            ) : filtered.map((r, i) => (
              <div key={i} style={{
                padding: "6px 10px", marginBottom: 4,
                border: `1px solid ${P.lightGray}`, borderRadius: 2,
                display: "flex", gap: 8, alignItems: "flex-start",
              }}>
                <span style={{ fontSize: 14, flexShrink: 0, marginTop: 1 }}>{REACTION_ICONS[r.type].icon}</span>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 10, color: P.darkGray, marginBottom: 2 }}><strong>{r.reader}</strong></div>
                  <div style={{ fontSize: 11, lineHeight: 1.4 }}>{r.text}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

// ─── Invite Dialog ───
function InviteDialog({ onClose }) {
  const [email, setEmail] = useState("");
  const [note, setNote] = useState("I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag.");
  const [sent, setSent] = useState(false);

  return (
    <div style={{
      position: "fixed", inset: 0, zIndex: 10000,
      display: "flex", alignItems: "center", justifyContent: "center",
      background: "rgba(0,0,0,0.15)",
    }}>
      <div style={{
        width: 420, border: `2px solid ${P.black}`, background: P.white,
        borderRadius: 4, boxShadow: "4px 4px 0 rgba(0,0,0,0.3)",
      }}>
        {/* Title bar */}
        <div style={{
          padding: "3px 8px", borderBottom: `2px solid ${P.black}`,
          display: "flex", alignItems: "center", fontFamily: CHICAGO, fontSize: 12,
        }}>
          <div onClick={onClose} style={{
            width: 13, height: 13, border: `1.5px solid ${P.black}`, cursor: "pointer", marginRight: 8, borderRadius: 1,
          }} />
          <div style={{ flex: 1, height: 10, background: P.stripes, marginRight: 6 }} />
          <span>Invite Reader</span>
          <div style={{ flex: 1, height: 10, background: P.stripes, marginLeft: 6 }} />
        </div>
        <div style={{ padding: 16, fontFamily: GENEVA, fontSize: 12 }}>
          {!sent ? (
            <>
              <p style={{ margin: "0 0 12px", lineHeight: 1.5 }}>
                Send a personal invitation. Remember: you want readers who would
                genuinely read this topic, even if you didn't write it.
              </p>
              <label style={{ fontFamily: CHICAGO, fontSize: 10, display: "block", marginBottom: 4 }}>
                Email Address
              </label>
              <input value={email} onChange={(e) => setEmail(e.target.value)} placeholder="reader@example.com"
                style={{
                  width: "100%", boxSizing: "border-box", padding: "6px 8px",
                  border: `2px solid ${P.black}`, fontFamily: GENEVA, fontSize: 12,
                  background: P.white, marginBottom: 12, borderRadius: 1, outline: "none",
                }}
              />
              <label style={{ fontFamily: CHICAGO, fontSize: 10, display: "block", marginBottom: 4 }}>
                Personal Note
              </label>
              <textarea value={note} onChange={(e) => setNote(e.target.value)} rows={4}
                style={{
                  width: "100%", boxSizing: "border-box", padding: "6px 8px",
                  border: `2px solid ${P.black}`, fontFamily: GENEVA, fontSize: 11,
                  background: P.white, marginBottom: 12, resize: "vertical", borderRadius: 1,
                  lineHeight: 1.5, outline: "none",
                }}
              />
              <div style={{ display: "flex", justifyContent: "flex-end", gap: 8 }}>
                <MacButton onClick={onClose}>Cancel</MacButton>
                <MacButton primary disabled={!email.includes("@")} onClick={() => setSent(true)}>Send Invite</MacButton>
              </div>
            </>
          ) : (
            <div style={{ textAlign: "center", padding: "20px 0" }}>
              <div style={{ fontSize: 32, marginBottom: 8 }}>✉</div>
              <div style={{ fontFamily: CHICAGO, fontSize: 13, marginBottom: 6 }}>Invitation Sent!</div>
              <div style={{ fontSize: 11, color: P.darkGray, marginBottom: 16, lineHeight: 1.5 }}>
                Your reader will receive a link to review the draft.
                Remember, only about 1 in 4 invitees will engage — that's normal!
              </div>
              <MacButton onClick={onClose}>OK</MacButton>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

// ─── About dialog ───
function AboutDialog({ onClose }) {
  return (
    <div style={{
      position: "fixed", inset: 0, zIndex: 10000,
      display: "flex", alignItems: "center", justifyContent: "center",
      background: "rgba(0,0,0,0.15)",
    }}>
      <div style={{
        width: 340, border: `2px solid ${P.black}`, background: P.white, borderRadius: 4,
        boxShadow: "4px 4px 0 rgba(0,0,0,0.3)",
      }}>
        <div style={{
          padding: "3px 8px", borderBottom: `2px solid ${P.black}`,
          display: "flex", alignItems: "center", fontFamily: CHICAGO, fontSize: 12,
        }}>
          <div onClick={onClose} style={{
            width: 13, height: 13, border: `1.5px solid ${P.black}`, cursor: "pointer", marginRight: 8,
          }} />
          <div style={{ flex: 1, height: 10, background: P.stripes, marginRight: 6 }} />
          <span>About</span>
          <div style={{ flex: 1, height: 10, background: P.stripes, marginLeft: 6 }} />
        </div>
        <div style={{ padding: 24, textAlign: "center", fontFamily: GENEVA, fontSize: 12 }}>
          <div style={{
            width: 64, height: 64, border: `2px solid ${P.black}`, margin: "0 auto 12px",
            display: "flex", alignItems: "center", justifyContent: "center",
            fontFamily: CHICAGO, fontSize: 28, borderRadius: 4,
          }}>📝</div>
          <div style={{ fontFamily: CHICAGO, fontSize: 14, marginBottom: 4 }}>Draft Review</div>
          <div style={{ fontSize: 11, color: P.darkGray, marginBottom: 12 }}>Version 1.0</div>
          <div style={{ fontSize: 11, lineHeight: 1.6, marginBottom: 16 }}>
            A retro-styled article draft review system
            inspired by beta reading methodology.
            <br /><br />
            Collect structured feedback. Find draft-killers.
            Ship better writing.
          </div>
          <MacButton primary onClick={onClose}>OK</MacButton>
        </div>
      </div>
    </div>
  );
}

// ─── Main App ───
export default function App() {
  const [view, setView] = useState("dashboard"); // dashboard | article
  const [selectedArticle, setSelectedArticle] = useState(null);
  const [focusSection, setFocusSection] = useState(null);
  const [showInvite, setShowInvite] = useState(false);
  const [showAbout, setShowAbout] = useState(false);
  const [articles] = useState(ARTICLES);

  const handleSelectArticle = useCallback((id, sectionId = null) => {
    setSelectedArticle(articles.find((a) => a.id === id));
    setFocusSection(sectionId);
    setView("article");
  }, [articles]);

  return (
    <div style={{
      position: "fixed", inset: 0,
      background: `${P.bg} repeating-conic-gradient(${P.lightGray} 0% 25%, transparent 0% 50%) 0 0 / 4px 4px`,
      fontFamily: GENEVA, fontSize: 12, overflow: "hidden",
    }}>
      <MenuBar
        onNewArticle={() => {}}
        onAbout={() => setShowAbout(true)}
        currentView={view}
        setView={setView}
      />

      {view === "dashboard" && (
        <MacWindow title="Draft Review — Dashboard" maximized zIndex={2}>
          <Dashboard
            articles={articles}
            onSelectArticle={handleSelectArticle}
            onInvite={() => setShowInvite(true)}
          />
        </MacWindow>
      )}

      {view === "article" && selectedArticle && (
        <MacWindow title={`Review: ${selectedArticle.title}`} maximized zIndex={3}>
          <ArticleReader
            article={selectedArticle}
            focusSection={focusSection}
            onBack={() => { setView("dashboard"); setSelectedArticle(null); }}
          />
        </MacWindow>
      )}

      {showInvite && <InviteDialog onClose={() => setShowInvite(false)} />}
      {showAbout && <AboutDialog onClose={() => setShowAbout(false)} />}
    </div>
  );
}
