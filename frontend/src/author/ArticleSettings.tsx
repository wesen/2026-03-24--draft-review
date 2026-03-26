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
  onGenerateLink?: () => void;
  shareUrl?: string;
}

type AccessMode = "invite_link" | "link" | "password";
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
  onGenerateLink,
  shareUrl,
}: ArticleSettingsProps) {
  const [title, setTitle] = useState(article.title);
  const [intro, setIntro] = useState(article.intro || "");
  const [version, setVersion] = useState(article.version || "");
  const [status, setStatus] = useState<ArticleStatus>(article.status);
  const [accessMode, setAccessMode] = useState<AccessMode>("invite_link");
  const [seeReactions, setSeeReactions] = useState(true);
  const [seeNames, setSeeNames] = useState(false);
  const [showAuthorNote, setShowAuthorNote] = useState(true);
  const [enabledReactions, setEnabledReactions] = useState<
    Record<string, boolean>
  >(
    Object.fromEntries(REACTION_TYPES.map((r) => [r.type, true]))
  );
  const [requireNote, setRequireNote] = useState(false);
  const [allowAnon, setAllowAnon] = useState(true);
  const [copied, setCopied] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleCopy = () => {
    if (shareUrl) {
      navigator.clipboard.writeText(window.location.origin + shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

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

        {/* SHARING */}
        <div className="dr-settings__section-header">SHARING</div>

        <label className="dr-settings__label">Review Link</label>
        <div className="dr-settings__link-row">
          <input
            className="dr-input dr-settings__link-input"
            value={
              shareUrl
                ? `${window.location.origin}${shareUrl}`
                : "No link generated yet"
            }
            readOnly
          />
          <MacButton small onClick={handleCopy} disabled={!shareUrl}>
            {copied ? "Copied!" : "Copy"}
          </MacButton>
          <MacButton small onClick={onGenerateLink}>
            {shareUrl ? "Reset" : "Generate"}
          </MacButton>
        </div>

        <label className="dr-settings__label">
          Access
          <span className="dr-settings__coming-soon">coming soon</span>
        </label>
        <div className="dr-settings__radio-group">
          {(
            [
              { value: "invite_link", label: "Unique reader links (no login)" },
              { value: "link", label: "Anyone with link (no login)" },
              { value: "password", label: "Password protected" },
            ] as { value: AccessMode; label: string }[]
          ).map((opt) => (
            <label key={opt.value} className="dr-settings__radio">
              <input
                type="radio"
                name="access"
                checked={accessMode === opt.value}
                onChange={() => setAccessMode(opt.value)}
              />
              {opt.label}
            </label>
          ))}
        </div>

        <label className="dr-settings__label">
          Reader visibility
          <span className="dr-settings__coming-soon">coming soon</span>
        </label>
        <div className="dr-settings__checkbox-group">
          <label className="dr-settings__checkbox">
            <input
              type="checkbox"
              checked={seeReactions}
              onChange={(e) => setSeeReactions(e.target.checked)}
            />
            Readers can see each other's reactions
          </label>
          <label className="dr-settings__checkbox">
            <input
              type="checkbox"
              checked={seeNames}
              onChange={(e) => setSeeNames(e.target.checked)}
            />
            Readers can see each other's names
          </label>
          <label className="dr-settings__checkbox">
            <input
              type="checkbox"
              checked={showAuthorNote}
              onChange={(e) => setShowAuthorNote(e.target.checked)}
            />
            Show author note on welcome screen
          </label>
        </div>

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
