import { useState, useMemo } from "react";
import type { Article, Reaction, Reader, Section } from "../types";
import { MacButton } from "../chrome/MacButton";
import { ProgressBar } from "../primitives/ProgressBar";
import "./ArticleManager.css";

type SortKey = "recent" | "reactions" | "readers" | "status";

interface ArticleManagerProps {
  articles: Article[];
  readers: Reader[];
  reactions: Reaction[];
  onEdit: (id: string) => void;
  onSettings: (id: string) => void;
  onReview: (id: string) => void;
  onNewArticle: () => void;
  onInvite: (id: string) => void;
}

export function ArticleManager({
  articles,
  readers,
  reactions,
  onEdit,
  onSettings,
  onReview,
  onNewArticle,
  onInvite,
}: ArticleManagerProps) {
  const [search, setSearch] = useState("");
  const [sortKey, setSortKey] = useState<SortKey>("recent");

  const enriched = useMemo(() => {
    return articles.map((a) => {
      const articleReaders = readers.filter((r) => r.articleId === a.id);
      const articleReactions = reactions.filter((r) => r.articleId === a.id);
      const avgProgress =
        articleReaders.length > 0
          ? Math.round(
              articleReaders.reduce((acc, r) => acc + r.progress, 0) /
                articleReaders.length
            )
          : 0;

      // Draft-killer
      const { killer } = a.sections.reduce<{
        killer: Section | null;
        maxPain: number;
      }>(
        (acc, s) => {
          const pain = articleReactions.filter(
            (r) =>
              r.sectionId === s.id &&
              (r.type === "confused" || r.type === "slow")
          ).length;
          return pain > acc.maxPain ? { killer: s, maxPain: pain } : acc;
        },
        { killer: null, maxPain: 0 }
      );

      return {
        article: a,
        readerCount: articleReaders.length,
        reactionCount: articleReactions.length,
        avgProgress,
        bookKiller: killer,
      };
    });
  }, [articles, readers, reactions]);

  const filtered = useMemo(() => {
    let items = enriched;
    if (search) {
      const q = search.toLowerCase();
      items = items.filter((e) =>
        e.article.title.toLowerCase().includes(q)
      );
    }
    items = [...items].sort((a, b) => {
      switch (sortKey) {
        case "reactions":
          return b.reactionCount - a.reactionCount;
        case "readers":
          return b.readerCount - a.readerCount;
        case "status":
          return a.article.status.localeCompare(b.article.status);
        case "recent":
        default:
          return (
            new Date(b.article.updatedAt).getTime() -
            new Date(a.article.updatedAt).getTime()
          );
      }
    });
    return items;
  }, [enriched, search, sortKey]);

  const active = filtered.filter((e) => e.article.status !== "archived");
  const archived = filtered.filter((e) => e.article.status === "archived");

  return (
    <div className="dr-manager">
      {/* Toolbar */}
      <div className="dr-manager__toolbar">
        <MacButton onClick={onNewArticle}>+ New Article</MacButton>
        <div className="dr-manager__toolbar-spacer" />
        <input
          className="dr-input dr-manager__search"
          placeholder="Search\u2026"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
        <select
          className="dr-manager__sort"
          value={sortKey}
          onChange={(e) => setSortKey(e.target.value as SortKey)}
        >
          <option value="recent">Recent</option>
          <option value="reactions">Most Reactions</option>
          <option value="readers">Most Readers</option>
          <option value="status">Status</option>
        </select>
      </div>

      {/* Article list */}
      <div className="dr-manager__list">
        {active.length === 0 && (
          <div className="dr-manager__empty">
            <div style={{ fontSize: 32, marginBottom: 8 }}>{"\uD83D\uDCDD"}</div>
            <div className="dr-manager__empty-title">No articles yet</div>
            <div className="dr-manager__empty-text">
              Create your first article and start getting feedback.
            </div>
            <MacButton primary onClick={onNewArticle}>
              New Article {"\u2192"}
            </MacButton>
          </div>
        )}

        {active.map((e) => (
          <ArticleCard
            key={e.article.id}
            {...e}
            onEdit={() => onEdit(e.article.id)}
            onSettings={() => onSettings(e.article.id)}
            onReview={() => onReview(e.article.id)}
            onInvite={() => onInvite(e.article.id)}
          />
        ))}

        {archived.length > 0 && (
          <ArchivedSection items={archived} onReview={onReview} />
        )}
      </div>
    </div>
  );
}

function ArticleCard({
  article,
  readerCount,
  reactionCount,
  avgProgress,
  bookKiller,
  onEdit,
  onSettings,
  onReview,
  onInvite,
}: {
  article: Article;
  readerCount: number;
  reactionCount: number;
  avgProgress: number;
  bookKiller: Section | null;
  onEdit: () => void;
  onSettings: () => void;
  onReview: () => void;
  onInvite: () => void;
}) {
  const statusLabel = article.status.replace("_", " ");
  return (
    <div className="dr-manager__card">
      <div className="dr-manager__card-header">
        <span className="dr-manager__card-icon">{"\u270E"}</span>
        <span className="dr-manager__card-title">{article.title}</span>
        <span className="dr-manager__card-status">{statusLabel}</span>
      </div>
      <div className="dr-manager__card-meta">
        {article.sections.length} sections {"\u00B7"} {reactionCount} reactions{" "}
        {"\u00B7"} {readerCount} readers
      </div>
      <div className="dr-manager__card-meta">{article.version}</div>
      {readerCount > 0 ? (
        <div className="dr-manager__card-progress">
          <ProgressBar percent={avgProgress} height={10} width={200} />
          <span className="dr-manager__card-pct">{avgProgress}% avg read</span>
        </div>
      ) : (
        <div className="dr-manager__card-no-readers">
          No readers yet — invite your first beta reader!
        </div>
      )}
      {bookKiller && (
        <div className="dr-manager__card-killer">
          {"\u26A0"} Draft-killer in "{bookKiller.title}"
        </div>
      )}
      <div className="dr-manager__card-actions">
        <MacButton small onClick={onReview}>
          Review
        </MacButton>
        <MacButton small onClick={onEdit}>
          Edit
        </MacButton>
        <MacButton small onClick={onSettings}>
          Share
        </MacButton>
        {readerCount === 0 && (
          <MacButton small primary onClick={onInvite}>
            Invite {"\u2192"}
          </MacButton>
        )}
      </div>
    </div>
  );
}

function ArchivedSection({
  items,
  onReview,
}: {
  items: {
    article: Article;
    readerCount: number;
    reactionCount: number;
    avgProgress: number;
    bookKiller: Section | null;
  }[];
  onReview: (id: string) => void;
}) {
  const [expanded, setExpanded] = useState(false);

  return (
    <div className="dr-manager__archived">
      <div
        className="dr-manager__archived-header"
        onClick={() => setExpanded(!expanded)}
      >
        Archived ({items.length}) {expanded ? "\u25BC" : "\u25B6"}
      </div>
      {expanded &&
        items.map((e) => (
          <div key={e.article.id} className="dr-manager__archived-item">
            <span>{e.article.title}</span>
            <MacButton small onClick={() => onReview(e.article.id)}>
              View
            </MacButton>
          </div>
        ))}
    </div>
  );
}
