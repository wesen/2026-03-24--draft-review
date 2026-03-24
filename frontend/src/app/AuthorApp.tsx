import { useState, useCallback } from "react";
import { MacWindow, MenuBar } from "../chrome";
import { Dashboard, ArticleReader, InviteDialog } from "../author";
import { ReaderPage } from "../reader";
import {
  useGetArticlesQuery,
  useGetReadersQuery,
  useGetReactionsQuery,
} from "../api/articleApi";
import type { Article } from "../types";

type View = "dashboard" | "article" | "reader-preview";

const menus = [
  {
    label: "File",
    items: [
      { label: "New Article\u2026", shortcut: "\u2318N" },
      { label: "Import from Docs\u2026", shortcut: "\u2318I" },
      { divider: true, label: "" },
      { label: "Export Feedback\u2026", shortcut: "\u2318E" },
    ],
  },
  {
    label: "Edit",
    items: [
      { label: "Undo", shortcut: "\u2318Z" },
      { label: "Redo", shortcut: "\u2318\u21E7Z" },
    ],
  },
  { label: "View", items: [{ label: "Dashboard" }] },
  {
    label: "Help",
    items: [
      { label: "Beta Reading Guide" },
      { label: "Keyboard Shortcuts" },
    ],
  },
];

export function AuthorApp() {
  const [view, setView] = useState<View>("dashboard");
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [focusSection, setFocusSection] = useState<string | undefined>();
  const [showInvite, setShowInvite] = useState(false);

  const { data: articles = [] } = useGetArticlesQuery();

  // Use first article as fallback for initial load
  const activeArticleId =
    selectedArticle?.id || articles[0]?.id || "";
  const { data: readers = [] } = useGetReadersQuery(activeArticleId, {
    skip: !activeArticleId,
  });
  const { data: reactions = [] } = useGetReactionsQuery(activeArticleId, {
    skip: !activeArticleId,
  });

  const handleSelectArticle = useCallback(
    (id: string, sectionId?: string) => {
      const article = articles.find((a) => a.id === id);
      if (article) {
        setSelectedArticle(article);
        setFocusSection(sectionId);
        setView("article");
      }
    },
    [articles]
  );

  const totalReactions = reactions.length;

  return (
    <div className="dr-desktop">
      <MenuBar
        menus={menus}
        rightStatus={`${totalReactions} reaction${totalReactions !== 1 ? "s" : ""}`}
      />

      {view === "dashboard" && articles.length > 0 && (
        <MacWindow
          title="Draft Review \u2014 Dashboard"
          maximized
          zIndex={2}
        >
          <Dashboard
            articles={articles}
            readers={readers}
            reactions={reactions}
            onSelectArticle={handleSelectArticle}
            onInvite={() => setShowInvite(true)}
          />
        </MacWindow>
      )}

      {view === "article" && selectedArticle && (
        <MacWindow
          title={`Review: ${selectedArticle.title}`}
          maximized
          zIndex={3}
        >
          <ArticleReader
            article={selectedArticle}
            reactions={reactions.filter(
              (r) => r.articleId === selectedArticle.id
            )}
            focusSection={focusSection}
            onBack={() => {
              setView("dashboard");
              setSelectedArticle(null);
            }}
          />
        </MacWindow>
      )}

      {view === "reader-preview" && selectedArticle && (
        <MacWindow
          title={`${selectedArticle.title} \u2014 Reader View`}
          maximized
          zIndex={3}
        >
          <ReaderPage
            article={{
              id: selectedArticle.id,
              title: selectedArticle.title,
              author: selectedArticle.author,
              version: selectedArticle.version,
              intro: selectedArticle.intro,
              sections: selectedArticle.sections,
            }}
          />
        </MacWindow>
      )}

      {showInvite && (
        <InviteDialog onClose={() => setShowInvite(false)} />
      )}
    </div>
  );
}
