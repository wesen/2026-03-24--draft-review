import { useState } from "react";
import type { Article, Reaction } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import { Prose } from "../primitives/Prose";
import "./ArticleReader.css";

interface ArticleReaderProps {
  article: Article;
  reactions: Reaction[];
  focusSection?: string;
  onBack: () => void;
}

export function ArticleReader({
  article,
  reactions,
  focusSection,
  onBack,
}: ArticleReaderProps) {
  const [selectedSection, setSelectedSection] = useState(
    focusSection || article.sections[0]?.id
  );
  const [reactionFilter, setReactionFilter] = useState<string | null>(null);

  const section = article.sections.find((s) => s.id === selectedSection);
  const sectionReactions = reactions.filter(
    (r) => r.sectionId === selectedSection
  );
  const filtered = reactionFilter
    ? sectionReactions.filter((r) => r.type === reactionFilter)
    : sectionReactions;

  return (
    <div className="dr-article-reader">
      {/* Sidebar */}
      <div className="dr-article-reader__sidebar">
        <div className="dr-article-reader__sidebar-header">Sections</div>
        <div className="dr-article-reader__section-list">
          {article.sections.map((s, i) => {
            const sr = reactions.filter((r) => r.sectionId === s.id);
            const isActive = s.id === selectedSection;
            const painCount = sr.filter(
              (r) => r.type === "confused" || r.type === "slow"
            ).length;
            return (
              <div
                key={s.id}
                onClick={() => {
                  setSelectedSection(s.id);
                  setReactionFilter(null);
                }}
                className={`dr-article-reader__section-item ${
                  isActive ? "dr-article-reader__section-item--active" : ""
                }`}
              >
                <div className="dr-article-reader__section-title">
                  {i + 1}. {s.title}
                </div>
                <div className="dr-article-reader__section-counts">
                  {REACTION_TYPES.map((rt) => {
                    const c = sr.filter((r) => r.type === rt.type).length;
                    return c > 0 ? (
                      <span key={rt.type}>
                        {rt.icon}
                        {c}
                      </span>
                    ) : null;
                  })}
                  {painCount >= 2 && <span>{"\u26A0"}</span>}
                </div>
              </div>
            );
          })}
        </div>
        <div className="dr-article-reader__sidebar-footer">
          <MacButton onClick={onBack} small fullWidth>
            ← Dashboard
          </MacButton>
        </div>
      </div>

      {/* Main content */}
      <div className="dr-article-reader__main">
        {/* Article text */}
        <div className="dr-article-reader__text">
          <h2 className="dr-article-reader__heading">{section?.title}</h2>
          <div className="dr-article-reader__section-badge">
            Section{" "}
            {article.sections.findIndex((s) => s.id === selectedSection) + 1} of{" "}
            {article.sections.length}
          </div>
          {section?.paragraphs.map((p, i) => (
            <Prose key={i} className="dr-article-reader__paragraph">
              {p}
            </Prose>
          ))}
        </div>

        {/* Reactions panel */}
        <div className="dr-article-reader__reactions-panel">
          <div className="dr-article-reader__filter-tabs">
            <div
              onClick={() => setReactionFilter(null)}
              className={`dr-article-reader__filter-tab ${
                !reactionFilter
                  ? "dr-article-reader__filter-tab--active"
                  : ""
              }`}
            >
              All ({sectionReactions.length})
            </div>
            {REACTION_TYPES.map((rt) => {
              const c = sectionReactions.filter(
                (r) => r.type === rt.type
              ).length;
              return (
                <div
                  key={rt.type}
                  onClick={() =>
                    setReactionFilter(
                      reactionFilter === rt.type ? null : rt.type
                    )
                  }
                  className={`dr-article-reader__filter-tab ${
                    reactionFilter === rt.type
                      ? "dr-article-reader__filter-tab--active"
                      : ""
                  }`}
                >
                  {rt.icon} {rt.label} ({c})
                </div>
              );
            })}
          </div>
          <div className="dr-article-reader__reaction-list">
            {filtered.length === 0 ? (
              <div className="dr-article-reader__no-reactions">
                No {reactionFilter || ""} reactions for this section yet.
              </div>
            ) : (
              filtered.map((r, i) => {
                const rt = REACTION_TYPES.find((t) => t.type === r.type);
                return (
                  <div key={i} className="dr-article-reader__reaction-card">
                    <span className="dr-article-reader__reaction-icon">
                      {rt?.icon}
                    </span>
                    <div className="dr-article-reader__reaction-content">
                      <div className="dr-article-reader__reaction-author">
                        {r.readerName}
                      </div>
                      <div>{r.text}</div>
                    </div>
                  </div>
                );
              })
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
