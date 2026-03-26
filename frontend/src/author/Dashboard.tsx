import { useState } from "react";
import type { Article, Section, Reaction, Reader } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import { ProgressBar } from "../primitives/ProgressBar";
import "./Dashboard.css";

interface DashboardProps {
  articles: Article[];
  readers: Reader[];
  reactions: Reaction[];
  onSelectArticle: (id: string, sectionId?: string) => void;
  onEditArticle: (id: string) => void;
  onArticleSettings: (id: string) => void;
  onNewArticle: () => void;
  onViewArticles: () => void;
  onInvite: (id: string) => void;
}

export function Dashboard({
  articles,
  readers,
  reactions,
  onSelectArticle,
  onEditArticle,
  onArticleSettings,
  onNewArticle,
  onViewArticles,
  onInvite,
}: DashboardProps) {
  const [selectedId, setSelectedId] = useState(articles[0]?.id);
  const selected = articles.find((a) => a.id === selectedId) || articles[0];

  if (!selected) {
    return (
      <div className="dr-dashboard">
        <div className="dr-dashboard__empty-state">
          <div style={{ fontSize: 32, marginBottom: 8 }}>{"\u270E"}</div>
          <div className="dr-dashboard__empty-title">No articles yet</div>
          <div className="dr-dashboard__empty-text">
            Create your first article to start collecting feedback.
          </div>
          <MacButton primary onClick={onNewArticle}>
            New Article {"\u2192"}
          </MacButton>
        </div>
      </div>
    );
  }

  const articleReaders = readers.filter((r) => r.articleId === selected.id);
  const articleReactions = reactions.filter((r) => r.articleId === selected.id);
  const totalReactions = articleReactions.length;

  const avgProgress =
    articleReaders.length > 0
      ? Math.round(
          articleReaders.reduce((a, r) => a + r.progress, 0) /
            articleReaders.length
        )
      : 0;

  // Reactions per section
  const sectionReactions: Record<string, Reaction[]> = {};
  selected.sections.forEach((s) => {
    sectionReactions[s.id] = articleReactions.filter(
      (r) => r.sectionId === s.id
    );
  });

  // Draft-killer: section with most confused + slow
  const { bookKiller, maxPain } = selected.sections.reduce<{
    bookKiller: Section | null;
    maxPain: number;
  }>(
    (acc, s) => {
      const pain = (sectionReactions[s.id] || []).filter(
        (r) => r.type === "confused" || r.type === "slow"
      ).length;
      return pain > acc.maxPain ? { bookKiller: s, maxPain: pain } : acc;
    },
    { bookKiller: null, maxPain: 0 }
  );

  return (
    <div className="dr-dashboard dr-dashboard--split">
      {/* Article list sidebar */}
      <div className="dr-dashboard__sidebar">
        <div className="dr-dashboard__sidebar-header">
          <span>Articles</span>
          <MacButton small onClick={onViewArticles}>
            All
          </MacButton>
        </div>
        <div className="dr-dashboard__article-list">
          {articles.map((a) => {
            const isSelected = a.id === selectedId;
            const aReactions = reactions.filter((r) => r.articleId === a.id);
            const aReaders = readers.filter((r) => r.articleId === a.id);
            return (
              <div
                key={a.id}
                className={`dr-dashboard__article-row ${
                  isSelected ? "dr-dashboard__article-row--active" : ""
                }`}
                onClick={() => setSelectedId(a.id)}
              >
                <div className="dr-dashboard__article-info">
                  <span className="dr-dashboard__article-title">{a.title}</span>
                  <span className="dr-dashboard__article-meta">
                    <span className="dr-dashboard__article-status">
                      {a.status.replace("_", " ")}
                    </span>
                    {" \u00B7 "}
                    {aReaders.length}R {aReactions.length}F
                  </span>
                </div>
                <div className="dr-dashboard__article-actions">
                  <span
                    className="dr-dashboard__action-icon"
                    title="Edit"
                    onClick={(e) => { e.stopPropagation(); onEditArticle(a.id); }}
                  >
                    {"\u270E"}
                  </span>
                  <span
                    className="dr-dashboard__action-icon"
                    title="Settings"
                    onClick={(e) => { e.stopPropagation(); onArticleSettings(a.id); }}
                  >
                    {"\u2699"}
                  </span>
                  <span
                    className="dr-dashboard__action-icon"
                    title="Invite Reader"
                    onClick={(e) => { e.stopPropagation(); onInvite(a.id); }}
                  >
                    {"\u2709"}
                  </span>
                  <span
                    className="dr-dashboard__action-icon"
                    title="View Reviews"
                    onClick={(e) => { e.stopPropagation(); onSelectArticle(a.id); }}
                  >
                    {"\u25B6"}
                  </span>
                </div>
              </div>
            );
          })}
        </div>
        <div className="dr-dashboard__sidebar-footer">
          <MacButton fullWidth small onClick={onNewArticle}>
            + New Article
          </MacButton>
        </div>
      </div>

      {/* Main content for selected article */}
      <div className="dr-dashboard__main">
        {/* Stats bar */}
        <div className="dr-dashboard__stat-bar">
          <span>{"\u25C9"} {articleReaders.length} readers</span>
          <span>{"\u2726"} {totalReactions} reactions</span>
          <span>{"\u00A7"} {selected.sections.length} sections</span>
          <span>
            {"\u25B8"} {articleReaders.length ? `${avgProgress}%` : "\u2014"} avg
            progress
          </span>
          <span className="dr-dashboard__version-badge">{selected.version}</span>
        </div>

        {/* Draft-killer alert */}
        {bookKiller && maxPain > 0 && (
          <div className="dr-dashboard__alert">
            <div className="dr-dashboard__alert-title">
              {"\u26A0"} Potential Draft-Killer Detected
            </div>
            <div className="dr-dashboard__alert-body">
              <strong>&ldquo;{bookKiller.title}&rdquo;</strong> has the most confusion and
              pacing issues ({maxPain} flags). Focus your next revision here.
            </div>
            <MacButton
              onClick={() => onSelectArticle(selected.id, bookKiller!.id)}
              small
            >
              Review This Section {"\u2192"}
            </MacButton>
          </div>
        )}

        {/* Main grid: readers + reactions chart */}
        <div className="dr-dashboard__grid">
          {/* Readers panel */}
          <div className="dr-dashboard__panel">
            <div className="dr-dashboard__panel-header">
              <span>Readers</span>
            </div>
            <div className="dr-dashboard__panel-body">
              {articleReaders.length === 0 ? (
                <div className="dr-dashboard__empty">
                  <div style={{ fontSize: 28, marginBottom: 6 }}>{"\u2709"}</div>
                  No readers yet. Invite your first beta reader!
                </div>
              ) : (
                articleReaders.map((r) => (
                  <div key={r.id} className="dr-dashboard__reader-row">
                    <div className="dr-dashboard__avatar">{r.avatar}</div>
                    <span className="dr-dashboard__reader-name">{r.name}</span>
                    <ProgressBar percent={r.progress} width={80} height={10} />
                    <span className="dr-dashboard__reader-pct">
                      {r.progress}%
                    </span>
                  </div>
                ))
              )}
            </div>
          </div>

          {/* Reactions by section */}
          <div className="dr-dashboard__panel">
            <div className="dr-dashboard__panel-header">
              Reactions by Section
            </div>
            <div className="dr-dashboard__panel-body">
              {selected.sections.map((s) => {
                const sr = sectionReactions[s.id] || [];
                const isKiller = bookKiller && s.id === bookKiller.id;
                return (
                  <div key={s.id} className="dr-dashboard__section-row">
                    <div className="dr-dashboard__section-label">
                      <span style={{ fontWeight: isKiller ? "bold" : "normal" }}>
                        {isKiller && "\u26A0 "}
                        {s.title}
                      </span>
                      <span className="dr-dashboard__section-count">
                        {sr.length} reaction{sr.length !== 1 ? "s" : ""}
                      </span>
                    </div>
                    <div className="dr-dashboard__section-bar">
                      {REACTION_TYPES.map((rt) => {
                        const c = sr.filter((r) => r.type === rt.type).length;
                        if (c === 0) return null;
                        return (
                          <div
                            key={rt.type}
                            className={`dr-dashboard__bar-segment dr-dashboard__bar-segment--${rt.type}`}
                            title={`${c} ${rt.label}`}
                            style={{ flex: c }}
                          />
                        );
                      })}
                      {sr.length === 0 && (
                        <div className="dr-dashboard__bar-empty" />
                      )}
                    </div>
                  </div>
                );
              })}
              <div className="dr-dashboard__legend">
                {REACTION_TYPES.map((r) => (
                  <span key={r.type}>
                    {r.icon} {r.label}
                  </span>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Recent feedback */}
        <div className="dr-dashboard__panel" style={{ marginBottom: 16 }}>
          <div className="dr-dashboard__panel-header">Recent Feedback</div>
          <div className="dr-dashboard__feedback-list">
            {articleReactions.length === 0 ? (
              <div className="dr-dashboard__empty">
                No feedback yet. Share your article to start collecting reactions.
              </div>
            ) : (
              articleReactions
                .slice(-6)
                .reverse()
                .map((r, i) => {
                  const rt = REACTION_TYPES.find((t) => t.type === r.type);
                  const section = selected.sections.find(
                    (s) => s.id === r.sectionId
                  );
                  let paraExcerpt = "";
                  if (section && r.paragraphId) {
                    const match = r.paragraphId.match(/-p(\d+)$/);
                    if (match) {
                      const pIdx = parseInt(match[1], 10);
                      const pText = section.paragraphs[pIdx];
                      if (pText) {
                        const plain = pText.replace(/[#*_`>\[\]()]/g, "").trim();
                        paraExcerpt = plain.length > 60
                          ? plain.slice(0, 60).trimEnd() + "\u2026"
                          : plain;
                      }
                    }
                  }
                  return (
                    <div key={i} className="dr-dashboard__feedback-item">
                      <span className="dr-dashboard__feedback-icon">
                        {rt?.icon}
                      </span>
                      <div>
                        <div className="dr-dashboard__feedback-meta">
                          <strong>{r.readerName}</strong> on{" "}
                          <em>{section?.title}</em>
                        </div>
                        {paraExcerpt && (
                          <div className="dr-dashboard__feedback-para">
                            &para; &ldquo;{paraExcerpt}&rdquo;
                          </div>
                        )}
                        <div className="dr-dashboard__feedback-text">
                          {r.text}
                        </div>
                      </div>
                    </div>
                  );
                })
            )}
          </div>
        </div>

        <div className="dr-dashboard__bottom-actions">
          <MacButton primary onClick={() => onSelectArticle(selected.id)}>
            Open Full Review {"\u2192"}
          </MacButton>
        </div>
      </div>
    </div>
  );
}
