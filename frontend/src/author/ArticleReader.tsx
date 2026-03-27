import { useState, useRef, useCallback } from "react";
import type { Article, Reaction } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import {
  deriveSectionBlocks,
  getBlockIndexFromParagraphId,
} from "../lib/markdownBlocks";
import { Prose } from "../primitives/Prose";
import "./ArticleReader.css";

interface ArticleReaderProps {
  article: Article;
  reactions: Reaction[];
  focusSection?: string;
  onBack: () => void;
}

/** Format a timestamp as relative time (e.g., "2h ago", "3d ago") */
function timeAgo(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const seconds = Math.floor((now - then) / 1000);
  if (seconds < 60) return "just now";
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  const months = Math.floor(days / 30);
  return `${months}mo ago`;
}

/** Truncate plain text for block excerpts */
function excerpt(plain: string, maxLen = 80): string {
  if (plain.length <= maxLen) return plain;
  return plain.slice(0, maxLen).trimEnd() + "…";
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
  const [hoveredParagraph, setHoveredParagraph] = useState<string | null>(null);
  const [filteredParagraph, setFilteredParagraph] = useState<string | null>(
    null
  );

  const textRef = useRef<HTMLDivElement>(null);

  const section = article.sections.find((s) => s.id === selectedSection);
  const sectionReactions = reactions.filter(
    (r) => r.sectionId === selectedSection
  );
  const blocks = section ? deriveSectionBlocks(section) : [];

  // Apply type filter
  let filtered = reactionFilter
    ? sectionReactions.filter((r) => r.type === reactionFilter)
    : sectionReactions;

  // Apply paragraph filter
  if (filteredParagraph) {
    filtered = filtered.filter((r) => r.paragraphId === filteredParagraph);
  }

  // Group reactions by paragraph for display
  const groupedByParagraph = new Map<string, Reaction[]>();
  for (const r of filtered) {
    const group = groupedByParagraph.get(r.paragraphId) || [];
    group.push(r);
    groupedByParagraph.set(r.paragraphId, group);
  }

  // Build paragraph text lookup for excerpts
  const paragraphText = new Map<string, string>();
  blocks.forEach((block) => {
    paragraphText.set(block.id, block.plainText);
  });

  // Scroll a paragraph into view when hovering a reaction card
  const scrollParagraphIntoView = useCallback(
    (paragraphId: string) => {
      if (!textRef.current) return;
      const el = textRef.current.querySelector(
        `[data-paragraph-id="${paragraphId}"]`
      );
      el?.scrollIntoView({ behavior: "smooth", block: "nearest" });
    },
    []
  );

  const handleParagraphClick = (paragraphId: string) => {
    setFilteredParagraph(
      filteredParagraph === paragraphId ? null : paragraphId
    );
  };

  const clearParagraphFilter = () => {
    setFilteredParagraph(null);
  };

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
                  setFilteredParagraph(null);
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
        <div className="dr-article-reader__text" ref={textRef}>
          <h2 className="dr-article-reader__heading">{section?.title}</h2>
          <div className="dr-article-reader__section-badge">
            Section{" "}
            {article.sections.findIndex((s) => s.id === selectedSection) + 1} of{" "}
            {article.sections.length}
          </div>
          {blocks.map((block) => {
            const pId = block.id;
            const pReactionCount = sectionReactions.filter(
              (r) => r.paragraphId === pId
            ).length;
            const isHovered = hoveredParagraph === pId;
            const isFiltered = filteredParagraph === pId;
            return (
              <div
                key={pId}
                data-paragraph-id={pId}
                className={`dr-review-para ${
                  isHovered ? "dr-review-para--active" : ""
                } ${pReactionCount > 0 ? "dr-review-para--has-reactions" : ""} ${
                  isFiltered ? "dr-review-para--filtered" : ""
                }`}
                onMouseEnter={() => setHoveredParagraph(pId)}
                onMouseLeave={() => setHoveredParagraph(null)}
                onClick={() => handleParagraphClick(pId)}
              >
                <Prose className="dr-article-reader__paragraph">
                  {block.markdown}
                </Prose>
                {pReactionCount > 0 && (
                  <span className="dr-review-para__badge">
                    {pReactionCount}
                  </span>
                )}
              </div>
            );
          })}
        </div>

        {/* Reactions panel */}
        <div className="dr-article-reader__reactions-panel">
          <div className="dr-article-reader__reactions-header">Reactions</div>
          <div className="dr-article-reader__filter-tabs">
            <div
              onClick={() => {
                setReactionFilter(null);
                setFilteredParagraph(null);
              }}
              className={`dr-article-reader__filter-tab ${
                !reactionFilter && !filteredParagraph
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
                  title={`${rt.label} (${c})`}
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
                  {rt.icon} {c}
                </div>
              );
            })}
          </div>

          {/* Paragraph filter indicator */}
          {filteredParagraph && (
            <div className="dr-article-reader__para-filter-bar">
              <span>
                Showing reactions for: &ldquo;
                {excerpt(paragraphText.get(filteredParagraph) || "")}
                &rdquo;
              </span>
              <span
                className="dr-article-reader__para-filter-clear"
                onClick={clearParagraphFilter}
              >
                Clear
              </span>
            </div>
          )}

          <div className="dr-article-reader__reaction-list">
            {filtered.length === 0 ? (
              <div className="dr-article-reader__no-reactions">
                No {reactionFilter || ""} reactions for this section yet.
              </div>
            ) : (
              Array.from(groupedByParagraph.entries()).map(
                ([pId, pReactions]) => (
                  <div key={pId} className="dr-article-reader__para-group">
                    <div className="dr-article-reader__para-excerpt">
                      &para; &ldquo;{excerpt(paragraphText.get(pId) || "")}
                      &rdquo;
                    </div>
                    {pReactions.map((r, i) => {
                      const rt = REACTION_TYPES.find(
                        (t) => t.type === r.type
                      );
                      const isCardHighlighted = hoveredParagraph === r.paragraphId;
                      const reactionBlockIndex = getBlockIndexFromParagraphId(r.paragraphId);
                      return (
                        <div
                          key={i}
                          className={`dr-article-reader__reaction-card ${
                            isCardHighlighted
                              ? "dr-article-reader__reaction-card--highlighted"
                              : ""
                          }`}
                          onMouseEnter={() => {
                            setHoveredParagraph(r.paragraphId);
                            scrollParagraphIntoView(r.paragraphId);
                          }}
                          onMouseLeave={() => setHoveredParagraph(null)}
                        >
                          <span className="dr-article-reader__reaction-icon">
                            {rt?.icon}
                          </span>
                          <div className="dr-article-reader__reaction-content">
                            <div className="dr-article-reader__reaction-meta">
                              <span className="dr-article-reader__reaction-author">
                                {r.readerName}
                              </span>
                              <span className="dr-article-reader__reaction-time">
                                {timeAgo(r.createdAt)}
                              </span>
                            </div>
                            {reactionBlockIndex !== null && (
                              <div className="dr-article-reader__reaction-meta">
                                Block {reactionBlockIndex + 1}
                              </div>
                            )}
                            <div>{r.text}</div>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                )
              )
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
