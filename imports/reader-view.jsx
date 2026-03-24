import { useState, useEffect, useRef, useCallback } from "react";

const CHICAGO = `"Chicago_12","ChicagoFLF","Geneva",monospace`;
const GENEVA = `"Geneva","Monaco",monospace`;

const P = {
  bg: "#e8e8e8",
  white: "#ffffff",
  black: "#000000",
  gray: "#a0a0a0",
  dark: "#555555",
  light: "#c0c0c0",
  stripes:
    "repeating-linear-gradient(0deg,#fff 0px,#fff 1px,#000 1px,#000 2px)",
};

// ── sample article ──
const ARTICLE = {
  title: "Why Design Systems Fail",
  author: "Alex Chen",
  version: "Draft 3 — March 2026",
  intro:
    "Thanks for reading an early draft! This is a work‑in‑progress that needs your honest reactions. As you read, highlight anything that feels useful, confusing, or slow — your feedback will directly shape the next revision.",
  sections: [
    {
      id: "s1",
      title: "Introduction",
      paragraphs: [
        "Design systems promise consistency and speed, but many teams abandon them within a year. This article explores why that happens and what you can do about it.",
        "The problem isn't usually technical — it's organizational. A design system is a product, and like any product, it fails when nobody owns it, nobody documents it, and nobody listens to its users.",
      ],
    },
    {
      id: "s2",
      title: "The Adoption Cliff",
      paragraphs: [
        "Most design systems see strong initial adoption. Engineers are excited, designers feel heard, and leadership is optimistic about the efficiency gains on the horizon.",
        "But around month six, usage plateaus. Teams start creating workarounds, one‑off components multiply, and the system starts to feel like a constraint rather than an enabler. The honeymoon is over.",
        "This is the adoption cliff, and it kills more design systems than bad component APIs ever will. The root cause is almost always a gap between what the system offers and what teams actually need day‑to‑day.",
      ],
    },
    {
      id: "s3",
      title: "Governance Gaps",
      paragraphs: [
        "Without clear ownership, design systems drift. Who decides when to add a new component? Who approves breaking changes? Who triages the 47 open GitHub issues?",
        "When these questions go unanswered, teams lose trust in the system and start going rogue. They fork components, inline styles, and build parallel mini‑systems inside their own repos.",
        "The fix isn't more rules — it's more clarity. A lightweight decision‑making framework, published and linked from your docs, can prevent months of confusion.",
      ],
    },
    {
      id: "s4",
      title: "The Documentation Problem",
      paragraphs: [
        "Even well‑built systems fail without good docs. If a developer can't figure out how to use a component in under two minutes, they'll build their own. Documentation isn't a nice‑to‑have — it's the product.",
        "The best design system docs don't just list props. They show real examples, explain when to use (and when not to use) each component, and link to the Figma source so designers and engineers are always looking at the same thing.",
      ],
    },
    {
      id: "s5",
      title: "Making It Work",
      paragraphs: [
        "The teams that succeed treat their design system like a product, not a project. They have dedicated maintainers, regular release cycles, and feedback loops with the teams who consume the system.",
        "Executive sponsorship matters too — not as a rubber stamp, but as air cover. When a product team wants to skip the system to hit a deadline, someone senior needs to weigh the long‑term cost.",
        "If you're starting a design system today, start small: one token file, one button, one clear owner. Grow it only when real demand pulls you forward. That patience is the hardest part — and the most important.",
      ],
    },
  ],
};

const REACTIONS = [
  { type: "useful", icon: "★", label: "Useful!", desc: "I can use this" },
  { type: "confused", icon: "?", label: "Confused", desc: "I don't follow" },
  { type: "slow", icon: "◎", label: "Slow", desc: "This drags" },
  { type: "favorite", icon: "♥", label: "Love it", desc: "Highlight" },
];

// ── chrome helpers ──
function TitleBar({ title, onClose, extra }) {
  return (
    <div
      style={{
        background: P.white,
        borderBottom: `2px solid ${P.black}`,
        padding: "2px 6px",
        display: "flex",
        alignItems: "center",
        minHeight: 22,
        userSelect: "none",
        flexShrink: 0,
      }}
    >
      {onClose && (
        <div
          onClick={onClose}
          style={{
            width: 14,
            height: 14,
            border: `1.5px solid ${P.black}`,
            borderRadius: 1,
            cursor: "pointer",
            marginRight: 8,
            flexShrink: 0,
          }}
        />
      )}
      <div style={{ flex: 1, height: 10, background: P.stripes, marginRight: 8 }} />
      <span
        style={{
          fontFamily: CHICAGO,
          fontSize: 12,
          fontWeight: "bold",
          whiteSpace: "nowrap",
          letterSpacing: 0.3,
        }}
      >
        {title}
      </span>
      <div style={{ flex: 1, height: 10, background: P.stripes, marginLeft: 8 }} />
      {extra}
    </div>
  );
}

function Btn({ children, onClick, primary, small, disabled, style = {} }) {
  return (
    <button
      disabled={disabled}
      onClick={onClick}
      style={{
        fontFamily: CHICAGO,
        fontSize: small ? 10 : 11,
        padding: small ? "2px 10px" : "4px 16px",
        border: `2px solid ${P.black}`,
        borderRadius: 8,
        background: primary ? P.black : P.white,
        color: primary ? P.white : P.black,
        cursor: disabled ? "default" : "pointer",
        opacity: disabled ? 0.4 : 1,
        boxShadow: primary ? "inset 0 0 0 1px #fff, 0 0 0 2px #000" : "none",
        ...style,
      }}
    >
      {children}
    </button>
  );
}

// ── inline reaction badge (already‑submitted) ──
function ReactionBadge({ type, text, onRemove }) {
  const r = REACTIONS.find((rx) => rx.type === type);
  return (
    <div
      style={{
        display: "inline-flex",
        alignItems: "flex-start",
        gap: 6,
        padding: "5px 10px",
        border: `1.5px solid ${P.black}`,
        borderRadius: 2,
        background: P.white,
        marginTop: 6,
        marginRight: 6,
        maxWidth: 340,
        fontSize: 11,
        fontFamily: GENEVA,
        lineHeight: 1.45,
      }}
    >
      <span style={{ fontSize: 13, flexShrink: 0, marginTop: -1 }}>{r?.icon}</span>
      <span style={{ flex: 1 }}>{text}</span>
      <span
        onClick={onRemove}
        style={{
          cursor: "pointer",
          fontFamily: CHICAGO,
          fontSize: 9,
          marginTop: 1,
          opacity: 0.5,
          flexShrink: 0,
        }}
      >
        ✕
      </span>
    </div>
  );
}

// ── paragraph component with reaction affordance ──
function Paragraph({ text, pid, reactions, onReact, onRemove }) {
  const [hovered, setHovered] = useState(false);
  const [picking, setPicking] = useState(false);
  const [chosenType, setChosenType] = useState(null);
  const [comment, setComment] = useState("");
  const inputRef = useRef(null);
  const myReactions = reactions.filter((r) => r.pid === pid);

  const submit = () => {
    if (!chosenType) return;
    onReact({ pid, type: chosenType, text: comment.trim() || REACTIONS.find((r) => r.type === chosenType)?.desc });
    setPicking(false);
    setChosenType(null);
    setComment("");
  };

  useEffect(() => {
    if (chosenType && inputRef.current) inputRef.current.focus();
  }, [chosenType]);

  return (
    <div
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => { setHovered(false); }}
      style={{ position: "relative", marginBottom: 6 }}
    >
      {/* The text */}
      <p
        style={{
          margin: 0,
          fontSize: 13,
          lineHeight: 1.85,
          color: P.black,
          padding: "4px 0 4px 38px",
          borderLeft: myReactions.length > 0 ? `3px solid ${P.black}` : "3px solid transparent",
          transition: "border-color .15s",
        }}
      >
        {text}
      </p>

      {/* + button on hover */}
      {hovered && !picking && (
        <div
          onClick={() => setPicking(true)}
          style={{
            position: "absolute",
            left: 6,
            top: 4,
            width: 22,
            height: 22,
            border: `1.5px solid ${P.black}`,
            borderRadius: "50%",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer",
            background: P.white,
            fontFamily: CHICAGO,
            fontSize: 13,
            lineHeight: 1,
            boxShadow: "1px 1px 0 rgba(0,0,0,.15)",
          }}
        >
          +
        </div>
      )}

      {/* Reaction picker */}
      {picking && (
        <div
          style={{
            marginLeft: 38,
            marginTop: 4,
            border: `2px solid ${P.black}`,
            borderRadius: 3,
            background: P.white,
            padding: 10,
            boxShadow: "3px 3px 0 rgba(0,0,0,.2)",
          }}
        >
          {/* type selector */}
          <div style={{ display: "flex", gap: 4, marginBottom: 8, flexWrap: "wrap" }}>
            {REACTIONS.map((r) => (
              <div
                key={r.type}
                onClick={() => setChosenType(r.type)}
                style={{
                  padding: "4px 12px",
                  border: `1.5px solid ${P.black}`,
                  borderRadius: 2,
                  cursor: "pointer",
                  background: chosenType === r.type ? P.black : P.white,
                  color: chosenType === r.type ? P.white : P.black,
                  fontFamily: CHICAGO,
                  fontSize: 10,
                  display: "flex",
                  gap: 5,
                  alignItems: "center",
                  userSelect: "none",
                }}
              >
                <span style={{ fontSize: 13 }}>{r.icon}</span>
                {r.label}
              </div>
            ))}
          </div>
          {/* comment */}
          {chosenType && (
            <>
              <textarea
                ref={inputRef}
                value={comment}
                onChange={(e) => setComment(e.target.value)}
                onKeyDown={(e) => { if (e.key === "Enter" && !e.shiftKey) { e.preventDefault(); submit(); } }}
                placeholder="Add a note (optional) — what specifically?"
                rows={2}
                style={{
                  width: "100%",
                  boxSizing: "border-box",
                  fontFamily: GENEVA,
                  fontSize: 11,
                  border: `1.5px solid ${P.black}`,
                  padding: "5px 8px",
                  resize: "none",
                  borderRadius: 1,
                  lineHeight: 1.5,
                  outline: "none",
                }}
              />
              <div style={{ display: "flex", justifyContent: "flex-end", gap: 6, marginTop: 6 }}>
                <Btn small onClick={() => { setPicking(false); setChosenType(null); setComment(""); }}>
                  Cancel
                </Btn>
                <Btn small primary onClick={submit}>
                  Submit
                </Btn>
              </div>
            </>
          )}
        </div>
      )}

      {/* Existing reactions */}
      {myReactions.length > 0 && (
        <div style={{ paddingLeft: 38 }}>
          {myReactions.map((r, i) => (
            <ReactionBadge key={i} type={r.type} text={r.text} onRemove={() => onRemove(r)} />
          ))}
        </div>
      )}
    </div>
  );
}

// ── section navigation dots ──
function SectionNav({ sections, current, readSections, onPick }) {
  return (
    <div style={{ display: "flex", gap: 3, alignItems: "center" }}>
      {sections.map((s, i) => {
        const isCurrent = s.id === current;
        const isRead = readSections.has(s.id);
        return (
          <div
            key={s.id}
            title={s.title}
            onClick={() => onPick(s.id)}
            style={{
              width: isCurrent ? 18 : 10,
              height: 10,
              border: `1.5px solid ${P.black}`,
              borderRadius: isCurrent ? 3 : "50%",
              background: isRead ? P.black : P.white,
              cursor: "pointer",
              transition: "all .15s",
            }}
          />
        );
      })}
    </div>
  );
}

// ── progress bar (retro hatched) ──
function HatchBar({ percent, w = "100%", h = 12 }) {
  return (
    <div
      style={{
        width: w,
        height: h,
        border: `1.5px solid ${P.black}`,
        background: P.white,
        position: "relative",
        borderRadius: 1,
        overflow: "hidden",
      }}
    >
      <div
        style={{
          position: "absolute",
          inset: "0 auto 0 0",
          width: `${percent}%`,
          background: `repeating-linear-gradient(45deg,${P.black} 0px,${P.black} 2px,${P.white} 2px,${P.white} 4px)`,
          transition: "width .3s",
        }}
      />
    </div>
  );
}

// ── completion dialog ──
function DoneDialog({ stats, onClose }) {
  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 10000,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "rgba(0,0,0,.15)",
      }}
    >
      <div
        style={{
          width: 380,
          border: `2px solid ${P.black}`,
          background: P.white,
          borderRadius: 4,
          boxShadow: "4px 4px 0 rgba(0,0,0,.3)",
          overflow: "hidden",
        }}
      >
        <TitleBar title="Review Complete" onClose={onClose} />
        <div style={{ padding: 24, textAlign: "center", fontFamily: GENEVA, fontSize: 12 }}>
          <div style={{ fontSize: 40, marginBottom: 8 }}>🎉</div>
          <div style={{ fontFamily: CHICAGO, fontSize: 14, marginBottom: 12 }}>
            Thank you for reading!
          </div>
          <div
            style={{
              display: "grid",
              gridTemplateColumns: "1fr 1fr",
              gap: 8,
              textAlign: "center",
              marginBottom: 16,
            }}
          >
            {REACTIONS.map((r) => (
              <div
                key={r.type}
                style={{
                  border: `1.5px solid ${P.black}`,
                  padding: "8px 6px",
                  borderRadius: 2,
                }}
              >
                <div style={{ fontSize: 18 }}>{r.icon}</div>
                <div style={{ fontFamily: CHICAGO, fontSize: 16, margin: "2px 0" }}>
                  {stats[r.type] || 0}
                </div>
                <div style={{ fontSize: 9, color: P.dark, textTransform: "uppercase", letterSpacing: 0.8 }}>
                  {r.label}
                </div>
              </div>
            ))}
          </div>
          <p style={{ fontSize: 11, color: P.dark, lineHeight: 1.6, margin: "0 0 16px" }}>
            Your reactions have been saved. The author will use
            them to improve the next revision. Every piece of
            feedback makes the final article better.
          </p>
          <Btn primary onClick={onClose}>
            Done
          </Btn>
        </div>
      </div>
    </div>
  );
}

// ── welcome splash ──
function WelcomeSplash({ article, onStart }) {
  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        height: "100%",
        padding: 32,
      }}
    >
      <div style={{ maxWidth: 460, textAlign: "center", fontFamily: GENEVA, fontSize: 12 }}>
        <div
          style={{
            width: 72,
            height: 72,
            border: `2px solid ${P.black}`,
            borderRadius: 6,
            margin: "0 auto 16px",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            fontSize: 36,
          }}
        >
          📝
        </div>
        <div style={{ fontFamily: CHICAGO, fontSize: 18, marginBottom: 4, letterSpacing: -0.3 }}>
          {article.title}
        </div>
        <div style={{ fontSize: 11, color: P.dark, marginBottom: 4 }}>
          by {article.author}
        </div>
        <div
          style={{
            display: "inline-block",
            padding: "2px 10px",
            border: `1px solid ${P.black}`,
            borderRadius: 2,
            fontSize: 9,
            fontFamily: CHICAGO,
            textTransform: "uppercase",
            letterSpacing: 1,
            marginBottom: 20,
          }}
        >
          {article.version}
        </div>

        <div
          style={{
            border: `2px solid ${P.black}`,
            borderRadius: 3,
            padding: 16,
            background: P.white,
            textAlign: "left",
            lineHeight: 1.7,
            marginBottom: 20,
          }}
        >
          <div style={{ fontFamily: CHICAGO, fontSize: 11, marginBottom: 8 }}>
            A note from the author:
          </div>
          {article.intro}
        </div>

        {/* How‑to */}
        <div
          style={{
            border: `2px solid ${P.black}`,
            borderRadius: 3,
            padding: 16,
            textAlign: "left",
            marginBottom: 24,
            background: `repeating-linear-gradient(45deg,${P.white} 0px,${P.white} 8px,#f7f7f7 8px,#f7f7f7 16px)`,
          }}
        >
          <div style={{ fontFamily: CHICAGO, fontSize: 11, marginBottom: 8 }}>
            How to leave feedback:
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "auto 1fr", gap: "6px 10px", fontSize: 11, lineHeight: 1.5 }}>
            {REACTIONS.map((r) => (
              <React.Fragment key={r.type}>
                <span style={{ fontSize: 15, textAlign: "center" }}>{r.icon}</span>
                <span>
                  <strong>{r.label}</strong> — {r.desc}
                </span>
              </React.Fragment>
            ))}
          </div>
          <div style={{ marginTop: 10, fontSize: 11, color: P.dark, lineHeight: 1.5 }}>
            Hover any paragraph and click <strong>+</strong> to react. 
            A short note explaining <em>why</em> is even more helpful. 
            Don't worry about typos or grammar — focus on what's useful, confusing, or slow.
          </div>
        </div>

        <Btn primary onClick={onStart}>
          Begin Reading →
        </Btn>
        <div style={{ fontSize: 10, color: P.gray, marginTop: 10 }}>
          {article.sections.length} sections · ~{article.sections.reduce((a, s) => a + s.paragraphs.length, 0) * 40} words · {Math.ceil(article.sections.reduce((a, s) => a + s.paragraphs.length, 0) * 0.4)} min read
        </div>
      </div>
    </div>
  );
}

// ── main reader app ──
export default function ReaderView() {
  const [started, setStarted] = useState(false);
  const [currentSection, setCurrentSection] = useState(ARTICLE.sections[0].id);
  const [readSections, setReadSections] = useState(new Set());
  const [reactions, setReactions] = useState([]);
  const [showDone, setShowDone] = useState(false);
  const scrollRef = useRef(null);

  const sectionIndex = ARTICLE.sections.findIndex((s) => s.id === currentSection);
  const section = ARTICLE.sections[sectionIndex];
  const isLast = sectionIndex === ARTICLE.sections.length - 1;
  const isFirst = sectionIndex === 0;
  const progress = ((readSections.size) / ARTICLE.sections.length) * 100;

  const goTo = useCallback(
    (id) => {
      setCurrentSection(id);
      if (scrollRef.current) scrollRef.current.scrollTop = 0;
    },
    []
  );

  const markRead = useCallback(() => {
    setReadSections((prev) => {
      const n = new Set(prev);
      n.add(currentSection);
      return n;
    });
  }, [currentSection]);

  const goNext = () => {
    markRead();
    if (!isLast) goTo(ARTICLE.sections[sectionIndex + 1].id);
  };
  const goPrev = () => {
    if (!isFirst) goTo(ARTICLE.sections[sectionIndex - 1].id);
  };

  const handleReact = (r) => {
    setReactions((prev) => [...prev, { ...r, section: currentSection, ts: Date.now() }]);
  };
  const handleRemove = (r) => {
    setReactions((prev) => prev.filter((x) => x !== r));
  };

  const finish = () => {
    markRead();
    setShowDone(true);
  };

  // stats for done dialog
  const stats = {};
  REACTIONS.forEach((r) => {
    stats[r.type] = reactions.filter((rx) => rx.type === r.type).length;
  });

  // ── keyboard nav ──
  useEffect(() => {
    const handler = (e) => {
      if (e.key === "ArrowRight" || e.key === "ArrowDown") {
        e.preventDefault();
        if (!isLast) goNext();
      }
      if (e.key === "ArrowLeft" || e.key === "ArrowUp") {
        e.preventDefault();
        if (!isFirst) goPrev();
      }
    };
    if (started) window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  });

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        background: `${P.bg} repeating-conic-gradient(${P.light} 0% 25%, transparent 0% 50%) 0 0 / 4px 4px`,
        fontFamily: GENEVA,
        fontSize: 12,
        display: "flex",
        flexDirection: "column",
      }}
    >
      {/* ── Menu bar ── */}
      <div
        style={{
          height: 22,
          background: P.white,
          borderBottom: `2px solid ${P.black}`,
          display: "flex",
          alignItems: "center",
          padding: "0 10px",
          fontFamily: CHICAGO,
          fontSize: 12,
          flexShrink: 0,
          gap: 14,
        }}
      >
        <span style={{ fontWeight: "bold", fontSize: 14 }}></span>
        <span>File</span>
        <span>Edit</span>
        <span>View</span>
        <span style={{ flex: 1 }} />
        <span style={{ fontSize: 10, color: P.dark }}>
          {reactions.length} reaction{reactions.length !== 1 ? "s" : ""} left
        </span>
      </div>

      {/* ── Window ── */}
      <div
        style={{
          flex: 1,
          margin: 0,
          display: "flex",
          flexDirection: "column",
          border: `2px solid ${P.black}`,
          borderTop: "none",
          background: P.white,
          overflow: "hidden",
        }}
      >
        <TitleBar title={`${ARTICLE.title} — Reader View`} />

        {!started ? (
          <div style={{ flex: 1, overflow: "auto" }}>
            <WelcomeSplash article={ARTICLE} onStart={() => setStarted(true)} />
          </div>
        ) : (
          <>
            {/* ── toolbar ── */}
            <div
              style={{
                padding: "6px 14px",
                borderBottom: `1.5px solid ${P.black}`,
                display: "flex",
                alignItems: "center",
                gap: 12,
                flexShrink: 0,
                flexWrap: "wrap",
              }}
            >
              <SectionNav
                sections={ARTICLE.sections}
                current={currentSection}
                readSections={readSections}
                onPick={(id) => goTo(id)}
              />
              <div style={{ flex: 1 }}>
                <HatchBar percent={progress} h={8} />
              </div>
              <span style={{ fontSize: 10, color: P.dark, fontFamily: CHICAGO }}>
                {Math.round(progress)}%
              </span>
            </div>

            {/* ── content ── */}
            <div ref={scrollRef} style={{ flex: 1, overflow: "auto" }}>
              <div style={{ maxWidth: 600, margin: "0 auto", padding: "28px 32px 60px" }}>
                {/* Section header */}
                <div
                  style={{
                    display: "flex",
                    alignItems: "baseline",
                    gap: 10,
                    marginBottom: 4,
                  }}
                >
                  <span
                    style={{
                      fontFamily: CHICAGO,
                      fontSize: 10,
                      color: P.dark,
                      border: `1px solid ${P.light}`,
                      borderRadius: 2,
                      padding: "1px 7px",
                      flexShrink: 0,
                    }}
                  >
                    {sectionIndex + 1} / {ARTICLE.sections.length}
                  </span>
                  <h2
                    style={{
                      fontFamily: CHICAGO,
                      fontSize: 18,
                      margin: 0,
                      letterSpacing: -0.3,
                    }}
                  >
                    {section.title}
                  </h2>
                </div>
                <div
                  style={{
                    height: 2,
                    background: P.stripes,
                    marginBottom: 22,
                    marginTop: 8,
                  }}
                />

                {/* Paragraphs */}
                {section.paragraphs.map((text, pi) => (
                  <Paragraph
                    key={`${section.id}-${pi}`}
                    text={text}
                    pid={`${section.id}-p${pi}`}
                    reactions={reactions.filter((r) => r.section === currentSection)}
                    onReact={handleReact}
                    onRemove={handleRemove}
                  />
                ))}

                {/* Section‑level summary of reactions */}
                {reactions.filter((r) => r.section === currentSection).length > 0 && (
                  <div
                    style={{
                      marginTop: 20,
                      padding: "8px 12px",
                      border: `1px solid ${P.light}`,
                      borderRadius: 2,
                      fontSize: 10,
                      color: P.dark,
                      display: "flex",
                      gap: 12,
                    }}
                  >
                    <span style={{ fontFamily: CHICAGO }}>This section:</span>
                    {REACTIONS.map((r) => {
                      const c = reactions.filter(
                        (rx) => rx.section === currentSection && rx.type === r.type
                      ).length;
                      return c > 0 ? (
                        <span key={r.type}>
                          {r.icon} {c}
                        </span>
                      ) : null;
                    })}
                  </div>
                )}

                {/* Navigation */}
                <div
                  style={{
                    marginTop: 32,
                    paddingTop: 16,
                    borderTop: `2px solid ${P.black}`,
                    display: "flex",
                    justifyContent: "space-between",
                    alignItems: "center",
                  }}
                >
                  <Btn onClick={goPrev} disabled={isFirst}>
                    ← Previous
                  </Btn>

                  {isLast ? (
                    <Btn primary onClick={finish}>
                      Finish Review ✓
                    </Btn>
                  ) : (
                    <Btn primary onClick={goNext}>
                      Next Section →
                    </Btn>
                  )}
                </div>

                {/* Keyboard hint */}
                <div
                  style={{
                    textAlign: "center",
                    marginTop: 12,
                    fontSize: 10,
                    color: P.gray,
                  }}
                >
                  Use arrow keys to navigate · hover paragraphs to react
                </div>
              </div>
            </div>
          </>
        )}
      </div>

      {showDone && <DoneDialog stats={stats} onClose={() => setShowDone(false)} />}
    </div>
  );
}
