import { useState, useCallback } from "react";
import { MacWindow, MenuBar } from "../chrome";
import {
  Dashboard,
  ArticleReader,
  ArticleEditor,
  ArticleSettings,
  ArticleManager,
  InviteDialog,
} from "../author";
import { ReaderPage } from "../reader";
import {
  useGetArticlesQuery,
  useGetReadersQuery,
  useGetReactionsQuery,
  useUpdateArticleMutation,
  useGenerateShareTokenMutation,
} from "../api/articleApi";
import type { Article } from "../types";

type View =
  | "dashboard"
  | "articles"
  | "article"
  | "edit"
  | "settings"
  | "reader-preview";

export function AuthorApp() {
  const [view, setView] = useState<View>("dashboard");
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [focusSection, setFocusSection] = useState<string | undefined>();
  const [showInvite, setShowInvite] = useState(false);
  const [shareUrl, setShareUrl] = useState<string | undefined>();

  const { data: articles = [] } = useGetArticlesQuery();
  const [updateArticle] = useUpdateArticleMutation();
  const [generateShareToken] = useGenerateShareTokenMutation();

  const activeArticleId = selectedArticle?.id || articles[0]?.id || "";
  const { data: readers = [] } = useGetReadersQuery(activeArticleId, {
    skip: !activeArticleId,
  });
  const { data: reactions = [] } = useGetReactionsQuery(activeArticleId, {
    skip: !activeArticleId,
  });

  const selectArticle = useCallback(
    (id: string) => {
      const a = articles.find((x) => x.id === id);
      if (a) setSelectedArticle(a);
    },
    [articles]
  );

  const handleSelectArticle = useCallback(
    (id: string, sectionId?: string) => {
      selectArticle(id);
      setFocusSection(sectionId);
      setView("article");
    },
    [selectArticle]
  );

  const goBack = () => {
    setView("dashboard");
    setSelectedArticle(null);
    setShareUrl(undefined);
  };

  const totalReactions = reactions.length;

  const menus = [
    {
      label: "File",
      items: [
        { label: "New Article\u2026", shortcut: "\u2318N" },
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
    {
      label: "View",
      items: [
        { label: "Dashboard", action: () => setView("dashboard") },
        { label: "Articles", action: () => setView("articles") },
      ],
    },
    {
      label: "Help",
      items: [
        { label: "Beta Reading Guide" },
        { label: "Keyboard Shortcuts" },
      ],
    },
  ];

  return (
    <div className="dr-desktop">
      <MenuBar
        menus={menus}
        rightStatus={`${totalReactions} reaction${totalReactions !== 1 ? "s" : ""}`}
      />

      {/* Dashboard */}
      {view === "dashboard" && articles.length > 0 && (
        <MacWindow title="Draft Review \u2014 Dashboard" maximized zIndex={2}>
          <Dashboard
            articles={articles}
            readers={readers}
            reactions={reactions}
            onSelectArticle={handleSelectArticle}
            onInvite={() => setShowInvite(true)}
          />
        </MacWindow>
      )}

      {/* Article Manager */}
      {view === "articles" && (
        <MacWindow title="My Articles" maximized zIndex={2}>
          <ArticleManager
            articles={articles}
            readers={readers}
            reactions={reactions}
            onEdit={(id) => {
              selectArticle(id);
              setView("edit");
            }}
            onSettings={(id) => {
              selectArticle(id);
              setShareUrl(undefined);
              setView("settings");
            }}
            onReview={(id) => handleSelectArticle(id)}
            onNewArticle={() => {}}
            onInvite={() => setShowInvite(true)}
          />
        </MacWindow>
      )}

      {/* Article Review (author view) */}
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
            onBack={goBack}
          />
        </MacWindow>
      )}

      {/* Article Editor */}
      {view === "edit" && selectedArticle && (
        <MacWindow
          title={`Edit: ${selectedArticle.title}`}
          maximized
          zIndex={3}
        >
          <ArticleEditor
            article={selectedArticle}
            onSave={async (updated) => {
              await updateArticle({
                id: updated.id,
                title: updated.title,
                sections: updated.sections,
                intro: updated.intro,
              });
              goBack();
            }}
            onBack={goBack}
            onPreview={(a) => {
              setSelectedArticle(a);
              setView("reader-preview");
            }}
          />
        </MacWindow>
      )}

      {/* Article Settings */}
      {view === "settings" && selectedArticle && (
        <MacWindow
          title={`Settings: ${selectedArticle.title}`}
          maximized
          zIndex={3}
        >
          <ArticleSettings
            article={selectedArticle}
            shareUrl={shareUrl}
            onSave={async (updates) => {
              await updateArticle({ id: selectedArticle.id, ...updates });
              goBack();
            }}
            onBack={goBack}
            onGenerateLink={async () => {
              const result = await generateShareToken(
                selectedArticle.id
              ).unwrap();
              setShareUrl(result.url);
            }}
          />
        </MacWindow>
      )}

      {/* Reader Preview */}
      {view === "reader-preview" && selectedArticle && (
        <MacWindow
          title={`${selectedArticle.title} \u2014 Reader View`}
          maximized
          zIndex={3}
          onClose={goBack}
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
