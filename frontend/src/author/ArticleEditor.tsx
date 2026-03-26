import { useState, useCallback, useEffect, useMemo, useRef } from "react";
import type { Article, Section } from "../types";
import { MacButton } from "../chrome/MacButton";
import "./ArticleEditor.css";

interface ArticleEditorProps {
  article: Article;
  onSave: (article: Article) => void;
  onBack: () => void;
  onPreview?: (article: Article) => void;
}

export function ArticleEditor({
  article: initialArticle,
  onSave,
  onBack,
  onPreview,
}: ArticleEditorProps) {
  const [article, setArticle] = useState<Article>(initialArticle);
  const [activeSectionId, setActiveSectionId] = useState(
    article.sections[0]?.id || ""
  );
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [showUnsavedConfirm, setShowUnsavedConfirm] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const activeSection = article.sections.find((s) => s.id === activeSectionId);

  // Serialize for dirty comparison — stable reference for initial state
  const initialJson = useMemo(
    () => JSON.stringify({ s: initialArticle.sections }),
    [initialArticle]
  );

  const parseParagraphs = useCallback((text: string) => {
    const paragraphs = text
      .split(/\n\n+/)
      .map((paragraph) => paragraph.trim())
      .filter((paragraph) => paragraph.length > 0);

    if (paragraphs.length === 0) {
      paragraphs.push("");
    }

    return paragraphs;
  }, []);

  const [contentDraft, setContentDraft] = useState(
    activeSection?.paragraphs.join("\n\n") || ""
  );

  const updateSection = useCallback(
    (sectionId: string, updates: Partial<Section>) => {
      setArticle((prev) => ({
        ...prev,
        sections: prev.sections.map((s) =>
          s.id === sectionId ? { ...s, ...updates } : s
        ),
      }));
    },
    []
  );

  useEffect(() => {
    setContentDraft(activeSection?.paragraphs.join("\n\n") || "");
  }, [activeSection]);

  // Auto-grow textarea
  useEffect(() => {
    const el = textareaRef.current;
    if (el) {
      el.style.height = "auto";
      el.style.height = Math.max(el.scrollHeight, 200) + "px";
    }
  }, [contentDraft]);

  const buildArticleWithCommittedDraft = useCallback(
    (source: Article) => {
      if (!activeSectionId) {
        return source;
      }

      return {
        ...source,
        sections: source.sections.map((section) =>
          section.id === activeSectionId
            ? { ...section, paragraphs: parseParagraphs(contentDraft) }
            : section
        ),
      };
    },
    [activeSectionId, contentDraft, parseParagraphs]
  );

  const commitContentDraft = useCallback(() => {
    setArticle((prev) => buildArticleWithCommittedDraft(prev));
  }, [buildArticleWithCommittedDraft]);

  // Dirty state: compare current article (with uncommitted draft) against initial
  const currentWithDraft = buildArticleWithCommittedDraft(article);
  const isDirty =
    JSON.stringify({ s: currentWithDraft.sections }) !== initialJson;

  const addSection = () => {
    const newSection: Section = {
      id: `s-new-${Date.now()}`,
      title: "New Section",
      paragraphs: [""],
    };
    setArticle((prev) => ({
      ...buildArticleWithCommittedDraft(prev),
      sections: [...buildArticleWithCommittedDraft(prev).sections, newSection],
    }));
    setActiveSectionId(newSection.id);
  };

  const deleteSection = () => {
    if (article.sections.length <= 1) return;
    setArticle((prev) => {
      const filtered = prev.sections.filter((s) => s.id !== activeSectionId);
      return { ...prev, sections: filtered };
    });
    setActiveSectionId(
      article.sections.find((s) => s.id !== activeSectionId)?.id || ""
    );
    setShowDeleteConfirm(false);
  };

  const moveSection = (fromIndex: number, direction: -1 | 1) => {
    const toIndex = fromIndex + direction;
    if (toIndex < 0 || toIndex >= article.sections.length) return;
    setArticle((prev) => {
      const committed = buildArticleWithCommittedDraft(prev);
      const sections = [...committed.sections];
      [sections[fromIndex], sections[toIndex]] = [
        sections[toIndex],
        sections[fromIndex],
      ];
      return { ...committed, sections };
    });
  };

  return (
    <div className="dr-editor">
      {/* Sidebar */}
      <div className="dr-editor__sidebar">
        <div className="dr-editor__sidebar-header">Sections</div>
        <div className="dr-editor__section-list">
          {article.sections.map((s, i) => (
            <div
              key={s.id}
              className={`dr-editor__section-item ${
                s.id === activeSectionId
                  ? "dr-editor__section-item--active"
                  : ""
              }`}
              onClick={() => {
                commitContentDraft();
                setActiveSectionId(s.id);
              }}
            >
              <div className="dr-editor__section-reorder">
                <span
                  className="dr-editor__reorder-btn"
                  onClick={(e) => {
                    e.stopPropagation();
                    moveSection(i, -1);
                  }}
                  title="Move up"
                >
                  {"\u25B2"}
                </span>
                <span
                  className="dr-editor__reorder-btn"
                  onClick={(e) => {
                    e.stopPropagation();
                    moveSection(i, 1);
                  }}
                  title="Move down"
                >
                  {"\u25BC"}
                </span>
              </div>
              <span className="dr-editor__section-label">
                {i + 1}. {s.title}
              </span>
            </div>
          ))}
        </div>
        <div className="dr-editor__sidebar-actions">
          <MacButton fullWidth onClick={addSection}>
            + Add Section
          </MacButton>
          <MacButton
            fullWidth
            disabled={article.sections.length <= 1}
            onClick={() => setShowDeleteConfirm(true)}
          >
            Delete Section
          </MacButton>
        </div>
        <div className="dr-editor__sidebar-footer">
          <MacButton
            fullWidth
            small
            onClick={() => {
              if (isDirty) {
                setShowUnsavedConfirm(true);
              } else {
                onBack();
              }
            }}
          >
            {"\u2190"} Back
          </MacButton>
        </div>
      </div>

      {/* Editor pane */}
      <div className="dr-editor__main">
        {activeSection ? (
          <div className="dr-editor__form">
            <label className="dr-editor__label">Title</label>
            <input
              className="dr-input dr-editor__title-input"
              value={activeSection.title}
              onChange={(e) =>
                updateSection(activeSection.id, { title: e.target.value })
              }
            />

            <label className="dr-editor__label">Content</label>
            <textarea
              ref={textareaRef}
              className="dr-textarea dr-editor__content-input"
              value={contentDraft}
              onBlur={commitContentDraft}
              onChange={(e) => setContentDraft(e.target.value)}
            />
            <div className="dr-editor__hint">
              Paragraph breaks (blank lines) become separate reactable
              paragraphs in the reader view.
            </div>

            <div className="dr-editor__actions">
              {isDirty && (
                <span className="dr-editor__dirty-indicator">
                  Unsaved changes
                </span>
              )}
              {onPreview && (
                <MacButton
                  onClick={() => {
                    const nextArticle = buildArticleWithCommittedDraft(article);
                    setArticle(nextArticle);
                    onPreview(nextArticle);
                  }}
                >
                  Preview as Reader
                </MacButton>
              )}
              <MacButton
                primary
                onClick={() => {
                  const nextArticle = buildArticleWithCommittedDraft(article);
                  setArticle(nextArticle);
                  onSave(nextArticle);
                }}
              >
                Save
              </MacButton>
            </div>
          </div>
        ) : (
          <div className="dr-editor__empty">
            No section selected. Add one to get started.
          </div>
        )}
      </div>

      {/* Delete confirmation */}
      {showDeleteConfirm && (
        <div className="dr-editor__confirm-overlay" onClick={() => setShowDeleteConfirm(false)}>
          <div className="dr-editor__confirm-dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dr-editor__confirm-title">Delete Section?</div>
            <div className="dr-editor__confirm-body">
              Are you sure you want to delete &ldquo;{activeSection?.title}&rdquo;?
              This cannot be undone.
            </div>
            <div className="dr-editor__confirm-actions">
              <MacButton onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </MacButton>
              <MacButton primary onClick={deleteSection}>
                Delete
              </MacButton>
            </div>
          </div>
        </div>
      )}

      {/* Unsaved changes confirmation */}
      {showUnsavedConfirm && (
        <div className="dr-editor__confirm-overlay" onClick={() => setShowUnsavedConfirm(false)}>
          <div className="dr-editor__confirm-dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dr-editor__confirm-title">Unsaved Changes</div>
            <div className="dr-editor__confirm-body">
              You have unsaved changes. Save before leaving?
            </div>
            <div className="dr-editor__confirm-actions">
              <MacButton onClick={() => { setShowUnsavedConfirm(false); onBack(); }}>
                Discard
              </MacButton>
              <MacButton
                primary
                onClick={() => {
                  const nextArticle = buildArticleWithCommittedDraft(article);
                  setArticle(nextArticle);
                  onSave(nextArticle);
                }}
              >
                Save
              </MacButton>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
