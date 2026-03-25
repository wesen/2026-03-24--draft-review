import { useState, useCallback } from "react";
import { MacWindow, MenuBar } from "../chrome";
import { MacButton } from "../chrome/MacButton";
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
  useCreateArticleMutation,
  useUpdateArticleMutation,
  useGenerateShareTokenMutation,
  useInviteReaderMutation,
} from "../api/articleApi";
import { useGetMeQuery } from "../api/authApi";
import { getBackendOrigin } from "../lib/backendOrigin";
import type { Article } from "../types";

type View =
  | "dashboard"
  | "articles"
  | "article"
  | "edit"
  | "settings"
  | "reader-preview";

export function AuthorApp() {
  const useMockApi = import.meta.env.VITE_USE_MSW === "1";
  const backendOrigin = getBackendOrigin();
  const [view, setView] = useState<View>("dashboard");
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [focusSection, setFocusSection] = useState<string | undefined>();
  const [showInvite, setShowInvite] = useState(false);
  const [shareUrl, setShareUrl] = useState<string | undefined>();

  const { data: meResponse, isLoading: isLoadingMe } = useGetMeQuery(undefined, {
    skip: useMockApi,
  });
  const me = meResponse?.data;
  const authReady = useMockApi || Boolean(me?.authenticated);

  const { data: articles = [] } = useGetArticlesQuery(undefined, {
    skip: !authReady,
  });
  const [createArticle] = useCreateArticleMutation();
  const [updateArticle] = useUpdateArticleMutation();
  const [generateShareToken] = useGenerateShareTokenMutation();
  const [inviteReader] = useInviteReaderMutation();

  const activeArticleId = selectedArticle?.id || articles[0]?.id || "";
  const { data: readers = [] } = useGetReadersQuery(activeArticleId, {
    skip: !authReady || !activeArticleId,
  });
  const { data: reactions = [] } = useGetReactionsQuery(activeArticleId, {
    skip: !authReady || !activeArticleId,
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

  const handleLogin = () => {
    window.location.assign(
      `${backendOrigin}/auth/login?return_to=${encodeURIComponent("/")}`
    );
  };

  const handleLogout = () => {
    window.location.assign(
      `${backendOrigin}/auth/logout?return_to=${encodeURIComponent("/")}`
    );
  };

  const handleNewArticle = useCallback(async () => {
    const article = await createArticle({
      title: "Untitled Article",
      author: me?.displayName || me?.preferredUsername || me?.email || "You",
      intro: "",
    }).unwrap();
    setSelectedArticle(article);
    setFocusSection(undefined);
    setView("edit");
  }, [createArticle, me?.displayName, me?.email, me?.preferredUsername]);

  const totalReactions = reactions.length;

  const menus = [
    {
      label: "File",
      items: [
        { label: "New Article\u2026", shortcut: "\u2318N", action: () => void handleNewArticle() },
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
    ...(!useMockApi && me?.authenticated
      ? [
          {
            label: "Account",
            items: [
              { label: `Signed in as ${me.displayName || me.email || "Author"}` },
              { divider: true, label: "" },
              { label: "Log Out", action: handleLogout },
            ],
          },
        ]
      : []),
  ];

  if (!useMockApi && isLoadingMe) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={[]} rightStatus="Checking session..." />
        <MacWindow title="Loading..." maximized>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              fontFamily: "var(--dr-font-body)",
              fontSize: "var(--dr-font-size-md)",
            }}
          >
            Loading your author session...
          </div>
        </MacWindow>
      </div>
    );
  }

  if (!useMockApi && !me?.authenticated) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={[]} rightStatus="Signed out" />
        <MacWindow title="Draft Review Login" maximized>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              flexDirection: "column",
              gap: 16,
              fontFamily: "var(--dr-font-body)",
              fontSize: "var(--dr-font-size-md)",
            }}
          >
            <div>Sign in to manage your draft reviews.</div>
            <MacButton primary onClick={handleLogin}>
              Sign In With Keycloak
            </MacButton>
          </div>
        </MacWindow>
      </div>
    );
  }

  return (
    <div className="dr-desktop">
      <MenuBar
        menus={menus}
        rightStatus={`${totalReactions} reaction${totalReactions !== 1 ? "s" : ""}${
          me?.displayName ? ` · ${me.displayName}` : ""
        }`}
      />

      {/* Dashboard */}
      {view === "dashboard" && (
        <MacWindow title="Draft Review \u2014 Dashboard" maximized zIndex={2}>
          <Dashboard
            articles={articles}
            readers={readers}
            reactions={reactions}
            onSelectArticle={handleSelectArticle}
            onEditArticle={(id) => {
              selectArticle(id);
              setView("edit");
            }}
            onArticleSettings={(id) => {
              selectArticle(id);
              setShareUrl(undefined);
              setView("settings");
            }}
            onViewArticles={() => setView("articles")}
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
            onNewArticle={() => void handleNewArticle()}
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
        <InviteDialog
          onClose={() => setShowInvite(false)}
          onInvite={async (email, note) => {
            if (!activeArticleId) {
              throw new Error("No active article selected");
            }
            const reader = await inviteReader({
              articleId: activeArticleId,
              email,
              note,
            }).unwrap();
            return {
              email: reader.email,
              inviteUrl: `${backendOrigin}/r/${reader.token}`,
            };
          }}
        />
      )}
    </div>
  );
}
