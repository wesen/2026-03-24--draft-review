import { useState } from "react";
import type { Article } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import "./ArticleSettings.css";

interface ArticleSettingsProps {
  article: Article;
  onSave: (updates: Partial<Article>) => void;
  onBack: () => void;
  onDelete?: () => void;
}

type ArticleStatus = "draft" | "in_review" | "complete" | "archived";

const STATUS_OPTIONS: { value: ArticleStatus; label: string }[] = [
  { value: "draft", label: "Draft" },
  { value: "in_review", label: "In Review" },
  { value: "complete", label: "Complete" },
  { value: "archived", label: "Archived" },
];

export function ArticleSettings({
  article,
  onSave,
  onBack,
  onDelete,
}: ArticleSettingsProps) {
  const [title, setTitle] = useState(article.title);
  const [intro, setIntro] = useState(article.intro || "");
  const [version, setVersion] = useState(article.version || "");
  const [status, setStatus] = useState<ArticleStatus>(article.status);
  const [enabledReactions, setEnabledReactions] = useState<
    Record<string, boolean>
  >(
    Object.fromEntries(REACTION_TYPES.map((r) => [r.type, true]))
  );
  const [requireNote, setRequireNote] = useState(false);
  const [allowAnon, setAllowAnon] = useState(true);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleSave = () => {
    onSave({ title, intro, version, status });
  };

  return (
    <div className="dr-settings">
      <div className="dr-settings__inner">
        {/* ARTICLE DETAILS */}
        <div className="dr-settings__section-header">ARTICLE DETAILS</div>

        <label className="dr-settings__label">Title</label>
        <input
          className="dr-input"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="Article title"
        />

        <label className="dr-settings__label">Author Intro</label>
        <textarea
          className="dr-textarea"
          value={intro}
          onChange={(e) => setIntro(e.target.value)}
          placeholder="A note for readers on the welcome screen"
          rows={3}
        />

        <label className="dr-settings__label">Version Label</label>
        <input
          className="dr-input"
          value={version}
          onChange={(e) => setVersion(e.target.value)}
          placeholder='e.g. "Draft 1", "Revision 2"'
        />

        {/* FEEDBACK */}
        <div className="dr-settings__section-header">FEEDBACK</div>

        <label className="dr-settings__label">
          Enabled reactions
          <span className="dr-settings__coming-soon">coming soon</span>
        </label>
        <div className="dr-settings__checkbox-group">
          {REACTION_TYPES.map((r) => (
            <label key={r.type} className="dr-settings__checkbox">
              <input
                type="checkbox"
                checked={enabledReactions[r.type]}
                onChange={(e) =>
                  setEnabledReactions((prev) => ({
                    ...prev,
                    [r.type]: e.target.checked,
                  }))
                }
              />
              {r.icon} {r.label}
            </label>
          ))}
        </div>

        <div className="dr-settings__checkbox-group">
          <label className="dr-settings__checkbox">
            <input
              type="checkbox"
              checked={requireNote}
              onChange={(e) => setRequireNote(e.target.checked)}
            />
            Require a note with each reaction
          </label>
          <label className="dr-settings__checkbox">
            <input
              type="checkbox"
              checked={allowAnon}
              onChange={(e) => setAllowAnon(e.target.checked)}
            />
            Allow anonymous reactions
          </label>
        </div>

        {/* STATUS */}
        <div className="dr-settings__section-header">STATUS</div>

        <div className="dr-settings__radio-group dr-settings__radio-group--inline">
          {STATUS_OPTIONS.map((opt) => (
            <label key={opt.value} className="dr-settings__radio">
              <input
                type="radio"
                name="status"
                checked={status === opt.value}
                onChange={() => setStatus(opt.value)}
              />
              {opt.label}
            </label>
          ))}
        </div>

        {/* DELETE */}
        {onDelete && (
          <>
            <div className="dr-settings__section-header">DELETE</div>
            <MacButton onClick={() => setShowDeleteConfirm(true)}>
              Delete Article
            </MacButton>
          </>
        )}

        {/* Actions */}
        <div className="dr-settings__actions">
          <MacButton onClick={onBack}>{"\u2190"} Back</MacButton>
          <MacButton primary onClick={handleSave}>
            Save Settings
          </MacButton>
        </div>
      </div>

      {/* Delete confirmation dialog */}
      {showDeleteConfirm && (
        <div className="dr-settings__confirm-overlay" onClick={() => setShowDeleteConfirm(false)}>
          <div className="dr-settings__confirm-dialog" onClick={(e) => e.stopPropagation()}>
            <div className="dr-settings__confirm-title">Delete Article?</div>
            <div className="dr-settings__confirm-body">
              Are you sure you want to delete &ldquo;{article.title}&rdquo;? This
              cannot be undone. All sections, readers, and reactions will be
              permanently removed.
            </div>
            <div className="dr-settings__confirm-actions">
              <MacButton onClick={() => setShowDeleteConfirm(false)}>
                Cancel
              </MacButton>
              <MacButton primary onClick={onDelete}>
                Delete
              </MacButton>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
