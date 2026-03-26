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
  useDeleteArticleMutation,
  useGenerateShareTokenMutation,
  useInviteReaderMutation,
} from "../api/articleApi";
import { useGetMeQuery } from "../api/authApi";
import { getBackendOrigin } from "../lib/backendOrigin";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  setView,
  selectArticle as selectArticleAction,
  setFocusSection,
  setPreviewArticle,
  goBack as goBackAction,
  openModal,
  closeModal,
} from "../store/uiSlice";

export function AuthorApp() {
  const useMockApi = import.meta.env.VITE_USE_MSW === "1";
  const backendOrigin = getBackendOrigin();
  const dispatch = useAppDispatch();

  // Redux state
  const view = useAppSelector((s) => s.ui.view);
  const selectedArticleId = useAppSelector((s) => s.ui.selectedArticleId);
  const focusSection = useAppSelector((s) => s.ui.focusSectionId);
  const previewArticle = useAppSelector((s) => s.ui.previewArticle);
  const activeModal = useAppSelector((s) => s.ui.activeModal);
  const showInvite = activeModal === "invite";

  // shareUrl is transient per-article — stays local
  const [shareUrl, setShareUrl] = useState<string | undefined>();

  const { data: meResponse, isLoading: isLoadingMe } = useGetMeQuery(undefined, {
    skip: useMockApi,
  });
  const me = meResponse?.data;
  const authReady = useMockApi || Boolean(me?.authenticated);

  const {
    data: articles = [],
    isLoading: isLoadingArticles,
    error: articlesError,
  } = useGetArticlesQuery(undefined, {
    skip: !authReady,
  });
  const [createArticle] = useCreateArticleMutation();
  const [updateArticle] = useUpdateArticleMutation();
  const [deleteArticle] = useDeleteArticleMutation();
  const [generateShareToken] = useGenerateShareTokenMutation();
  const [inviteReader] = useInviteReaderMutation();

  // Hydrate full Article from RTK Query cache
  const selectedArticle = articles.find((a) => a.id === selectedArticleId) || null;
  const activeArticleId = selectedArticleId || articles[0]?.id || "";

  const {
    data: readers = [],
    isLoading: isLoadingReaders,
  } = useGetReadersQuery(activeArticleId, {
    skip: !authReady || !activeArticleId,
  });
  const {
    data: reactions = [],
    isLoading: isLoadingReactions,
  } = useGetReactionsQuery(activeArticleId, {
    skip: !authReady || !activeArticleId,
  });

  const handleSelectArticle = useCallback(
    (id: string, sectionId?: string) => {
      dispatch(selectArticleAction(id));
      dispatch(setFocusSection(sectionId ?? null));
      dispatch(setView("article"));
    },
    [dispatch]
  );

  const goBack = useCallback(() => {
    dispatch(goBackAction());
    setShareUrl(undefined);
  }, [dispatch]);

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
    dispatch(selectArticleAction(article.id));
    dispatch(setFocusSection(null));
    dispatch(setView("edit"));
  }, [createArticle, dispatch, me?.displayName, me?.email, me?.preferredUsername]);

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
        { label: "Dashboard", action: () => dispatch(setView("dashboard")) },
        { label: "Articles", action: () => dispatch(setView("articles")) },
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

  if (isLoadingArticles) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={menus} rightStatus="Loading..." />
        <MacWindow title="Draft Review" maximized>
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
            Loading articles...
          </div>
        </MacWindow>
      </div>
    );
  }

  if (articlesError) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={menus} rightStatus="Error" />
        <MacWindow title={"Draft Review \u2014 Error"} maximized>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              flexDirection: "column",
              gap: 12,
              fontFamily: "var(--dr-font-body)",
              fontSize: "var(--dr-font-size-md)",
            }}
          >
            <div>Failed to load articles. Check your connection and try again.</div>
            <MacButton onClick={() => window.location.reload()}>
              Retry
            </MacButton>
          </div>
        </MacWindow>
      </div>
    );
  }

  const dataLoading = isLoadingReaders || isLoadingReactions;

  return (
    <div className="dr-desktop">
      <MenuBar
        menus={menus}
        rightStatus={`${
          dataLoading ? "Loading\u2026 · " : ""
        }${totalReactions} reaction${totalReactions !== 1 ? "s" : ""}${
          me?.displayName ? ` \u00B7 ${me.displayName}` : ""
        }`}
      />

      {/* Dashboard */}
      {view === "dashboard" && (
        <MacWindow title={"Draft Review \u2014 Dashboard"} maximized zIndex={2}>
          <Dashboard
            articles={articles}
            readers={readers}
            reactions={reactions}
            onSelectArticle={handleSelectArticle}
            onEditArticle={(id) => {
              dispatch(selectArticleAction(id));
              dispatch(setView("edit"));
            }}
            onArticleSettings={(id) => {
              dispatch(selectArticleAction(id));
              setShareUrl(undefined);
              dispatch(setView("settings"));
            }}
            onNewArticle={() => void handleNewArticle()}
            onViewArticles={() => dispatch(setView("articles"))}
            onInvite={() => dispatch(openModal("invite"))}
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
              dispatch(selectArticleAction(id));
              dispatch(setView("edit"));
            }}
            onSettings={(id) => {
              dispatch(selectArticleAction(id));
              setShareUrl(undefined);
              dispatch(setView("settings"));
            }}
            onReview={(id) => handleSelectArticle(id)}
            onNewArticle={() => void handleNewArticle()}
            onInvite={() => dispatch(openModal("invite"))}
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
            focusSection={focusSection ?? undefined}
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
              dispatch(setPreviewArticle(a));
              dispatch(setView("reader-preview"));
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
            onDelete={async () => {
              await deleteArticle(selectedArticle.id);
              goBack();
            }}
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
      {view === "reader-preview" && previewArticle && (
        <MacWindow
          title={`${previewArticle.title} \u2014 Reader View`}
          maximized
          zIndex={3}
          onClose={goBack}
        >
          <ReaderPage
            article={{
              id: previewArticle.id,
              title: previewArticle.title,
              author: previewArticle.author,
              version: previewArticle.version,
              intro: previewArticle.intro,
              sections: previewArticle.sections,
            }}
            readOnly
            onBackToEditor={() => dispatch(setView("edit"))}
          />
        </MacWindow>
      )}

      {showInvite && (
        <InviteDialog
          onClose={() => dispatch(closeModal())}
          shareUrl={shareUrl}
          onGenerateShareLink={async () => {
            if (!activeArticleId) {
              throw new Error("No active article selected");
            }
            const result = await generateShareToken(activeArticleId).unwrap();
            setShareUrl(result.url);
            return result.url;
          }}
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
