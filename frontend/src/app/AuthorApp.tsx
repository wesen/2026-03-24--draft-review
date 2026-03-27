import { useCallback } from "react";
import {
  useMatch,
  useNavigate,
  useSearchParams,
} from "react-router-dom";
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
  useUploadArticleAssetMutation,
} from "../api/articleApi";
import { useGetMeQuery } from "../api/authApi";
import { getBackendOrigin } from "../lib/backendOrigin";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  setPreviewArticle,
} from "../store/uiSlice";

export function AuthorApp() {
  const useMockApi = import.meta.env.VITE_USE_MSW === "1";
  const backendOrigin = getBackendOrigin();
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const previewArticle = useAppSelector((s) => s.ui.previewArticle);
  const routeSectionId = searchParams.get("section");

  const reviewMatch = useMatch("/articles/:articleId");
  const editMatch = useMatch("/articles/:articleId/edit");
  const settingsMatch = useMatch("/articles/:articleId/settings");
  const shareMatch = useMatch("/articles/:articleId/share");
  const previewMatch = useMatch("/articles/:articleId/preview");
  const articlesMatch = useMatch("/articles");

  const routeArticleId =
    shareMatch?.params.articleId ||
    previewMatch?.params.articleId ||
    settingsMatch?.params.articleId ||
    editMatch?.params.articleId ||
    reviewMatch?.params.articleId ||
    null;
  const currentView =
    shareMatch != null
      ? "share"
      : previewMatch != null
        ? "reader-preview"
        : settingsMatch != null
          ? "settings"
          : editMatch != null
            ? "edit"
            : reviewMatch != null
              ? "article"
              : articlesMatch != null
                ? "articles"
                : "dashboard";

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
  const [uploadArticleAsset] = useUploadArticleAssetMutation();

  const selectedArticle = articles.find((a) => a.id === routeArticleId) || null;
  const previewSourceArticle =
    previewArticle && previewArticle.id === routeArticleId
      ? previewArticle
      : selectedArticle;
  const activeArticleId = routeArticleId || articles[0]?.id || "";

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
      const query = sectionId
        ? `?section=${encodeURIComponent(sectionId)}`
        : "";
      navigate(`/articles/${id}${query}`);
    },
    [navigate]
  );

  const goBack = useCallback(() => {
    if (window.history.length > 1) {
      navigate(-1);
      return;
    }
    navigate("/");
  }, [navigate]);

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
    navigate(`/articles/${article.id}/edit`);
  }, [createArticle, me?.displayName, me?.email, me?.preferredUsername, navigate]);

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
        { label: "Dashboard", action: () => navigate("/") },
        { label: "Articles", action: () => navigate("/articles") },
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
      {currentView === "dashboard" && (
        <MacWindow title={"Draft Review \u2014 Dashboard"} maximized zIndex={2}>
          <Dashboard
            articles={articles}
            readers={readers}
            reactions={reactions}
            onSelectArticle={handleSelectArticle}
            onEditArticle={(id) => navigate(`/articles/${id}/edit`)}
            onArticleSettings={(id) => navigate(`/articles/${id}/settings`)}
            onNewArticle={() => void handleNewArticle()}
            onViewArticles={() => navigate("/articles")}
            onInvite={(id) => navigate(`/articles/${id}/share`)}
          />
        </MacWindow>
      )}

      {/* Article Manager */}
      {currentView === "articles" && (
        <MacWindow title="My Articles" maximized zIndex={2}>
          <ArticleManager
            articles={articles}
            readers={readers}
            reactions={reactions}
            onEdit={(id) => navigate(`/articles/${id}/edit`)}
            onSettings={(id) => navigate(`/articles/${id}/settings`)}
            onReview={(id) => handleSelectArticle(id)}
            onNewArticle={() => void handleNewArticle()}
            onInvite={(id) => navigate(`/articles/${id}/share`)}
          />
        </MacWindow>
      )}

      {/* Article Review (author view) */}
      {currentView === "article" && selectedArticle && (
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
            focusSection={routeSectionId ?? undefined}
            onBack={goBack}
          />
        </MacWindow>
      )}

      {/* Article Editor */}
      {currentView === "edit" && selectedArticle && (
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
              }).unwrap();
              navigate(`/articles/${updated.id}`);
            }}
            onBack={() => navigate(`/articles/${selectedArticle.id}`)}
            onUploadAsset={(file) =>
              uploadArticleAsset({
                articleId: selectedArticle.id,
                file,
              }).unwrap()
            }
            onPreview={(a) => {
              dispatch(setPreviewArticle(a));
              navigate(`/articles/${a.id}/preview`);
            }}
          />
        </MacWindow>
      )}

      {/* Article Settings */}
      {(currentView === "settings" || currentView === "share") && selectedArticle && (
        <MacWindow
          title={`Settings: ${selectedArticle.title}`}
          maximized
          zIndex={3}
        >
          <ArticleSettings
            article={selectedArticle}
            onSave={async (updates) => {
              await updateArticle({ id: selectedArticle.id, ...updates }).unwrap();
              navigate(`/articles/${selectedArticle.id}`);
            }}
            onBack={() => navigate(`/articles/${selectedArticle.id}`)}
            onDelete={async () => {
              await deleteArticle(selectedArticle.id).unwrap();
              navigate("/");
            }}
          />
        </MacWindow>
      )}

      {/* Reader Preview */}
      {currentView === "reader-preview" && previewSourceArticle && (
        <MacWindow
          title={`${previewSourceArticle.title} \u2014 Reader View`}
          maximized
          zIndex={3}
          onClose={goBack}
        >
          <ReaderPage
            article={{
              id: previewSourceArticle.id,
              title: previewSourceArticle.title,
              author: previewSourceArticle.author,
              version: previewSourceArticle.version,
              intro: previewSourceArticle.intro,
              sections: previewSourceArticle.sections,
            }}
            previewMode
            onBackToEditor={() => navigate(`/articles/${previewSourceArticle.id}/edit`)}
          />
        </MacWindow>
      )}

      {currentView === "share" && selectedArticle && (
        <InviteDialog
          shareUrl={
            selectedArticle.shareUrl
              ? `${backendOrigin}${selectedArticle.shareUrl}`
              : ""
          }
          onClose={() => {
            if (window.history.length > 1) {
              navigate(-1);
              return;
            }
            navigate(`/articles/${selectedArticle.id}/settings`);
          }}
          onGenerateShareLink={async () => {
            if (!activeArticleId) {
              throw new Error("No active article selected");
            }
            const result = await generateShareToken(activeArticleId).unwrap();
            return `${backendOrigin}${result.url}`;
          }}
          onInvite={async ({ identityMode, displayName, email, note }) => {
            if (!activeArticleId) {
              throw new Error("No active article selected");
            }
            const reader = await inviteReader({
              articleId: activeArticleId,
              identityMode,
              displayName,
              email,
              note,
            }).unwrap();
            return {
              name: reader.name,
              email: reader.email,
              identityMode: reader.identityMode || identityMode,
              inviteUrl: `${backendOrigin}/r/${reader.token}`,
            };
          }}
        />
      )}

      {routeArticleId && !selectedArticle && !isLoadingArticles && (
        <MacWindow title="Article Not Found" maximized zIndex={3}>
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
            <div>The requested article could not be found.</div>
            <MacButton onClick={() => navigate("/")}>Return to Dashboard</MacButton>
          </div>
        </MacWindow>
      )}
    </div>
  );
}
