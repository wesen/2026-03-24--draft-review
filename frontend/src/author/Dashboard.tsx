import { useState } from "react";
import type { Article, Section, Reaction, Reader } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import { ProgressBar } from "../primitives/ProgressBar";
import { StatCard } from "../primitives/StatCard";
import "./Dashboard.css";

interface DashboardProps {
  articles: Article[];
  readers: Reader[];
  reactions: Reaction[];
  onSelectArticle: (id: string, sectionId?: string) => void;
  onEditArticle: (id: string) => void;
  onArticleSettings: (id: string) => void;
  onViewArticles: () => void;
  onInvite: () => void;
}

export function Dashboard({
  articles,
  readers,
  reactions,
  onSelectArticle,
  onEditArticle,
  onArticleSettings,
  onViewArticles,
  onInvite,
}: DashboardProps) {
  const [selectedId, setSelectedId] = useState(articles[0]?.id);
  const selected = articles.find((a) => a.id === selectedId) || articles[0];

  if (!selected) {
    return (
      <div className="dr-dashboard">
        <div className="dr-dashboard__empty-state">
          <div style={{ fontSize: 32, marginBottom: 8 }}>{"\uD83D\uDCDD"}</div>
          <div className="dr-dashboard__empty-title">No articles yet</div>
          <div className="dr-dashboard__empty-text">
            Create your first article to start collecting feedback.
          </div>
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
    <div className="dr-dashboard">
      {/* Article selector + nav */}
      <div className="dr-dashboard__top-bar">
        <div className="dr-dashboard__tabs">
          {articles.map((a) => (
            <div
              key={a.id}
              onClick={() => setSelectedId(a.id)}
              className={`dr-dashboard__tab ${
                a.id === selectedId ? "dr-dashboard__tab--active" : ""
              }`}
            >
              <span className="dr-dashboard__tab-status">
                {a.status.replace("_", " ")}
              </span>
              {a.title}
            </div>
          ))}
        </div>
        <MacButton small onClick={onViewArticles}>
          All Articles
        </MacButton>
      </div>

      {/* Quick actions for selected article */}
      <div className="dr-dashboard__actions-bar">
        <MacButton small onClick={() => onEditArticle(selected.id)}>
          Edit
        </MacButton>
        <MacButton small onClick={() => onArticleSettings(selected.id)}>
          Share
        </MacButton>
        <MacButton small onClick={onInvite}>
          + Invite Reader
        </MacButton>
        <div style={{ flex: 1 }} />
        <span className="dr-dashboard__version-badge">{selected.version}</span>
      </div>

      {/* Stats row */}
      <div className="dr-dashboard__stats">
        <StatCard icon={"\u25C9"} value={articleReaders.length} label="Readers" />
        <StatCard icon={"\u2726"} value={totalReactions} label="Reactions" />
        <StatCard icon={"\u00A7"} value={selected.sections.length} label="Sections" />
        <StatCard
          icon={"\u25B8"}
          value={articleReaders.length ? `${avgProgress}%` : "\u2014"}
          label="Avg Progress"
        />
      </div>

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
                <div style={{ fontSize: 28, marginBottom: 6 }}>{"\uD83D\uDCED"}</div>
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

      {/* Draft-killer alert */}
      {bookKiller && maxPain > 0 && (
        <div className="dr-dashboard__alert">
          <div className="dr-dashboard__alert-title">
            {"\u26A0"} Potential Draft-Killer Detected
          </div>
          <div className="dr-dashboard__alert-body">
            <strong>"{bookKiller.title}"</strong> has the most confusion and
            pacing issues ({maxPain} flags). Readers are telling you this section
            needs work. Focus your next revision here.
          </div>
          <MacButton
            onClick={() => onSelectArticle(selected.id, bookKiller!.id)}
            small
          >
            Review This Section {"\u2192"}
          </MacButton>
        </div>
      )}

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
  );
}
