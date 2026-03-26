import { useEffect, useMemo, useState } from "react";
import { MacButton } from "../chrome/MacButton";
import type { InviteIdentityMode } from "../types";
import "./InviteDialog.css";

interface InviteResult {
  name: string;
  email?: string;
  identityMode: InviteIdentityMode;
  inviteUrl: string;
}

interface InviteDialogProps {
  onClose: () => void;
  shareUrl?: string;
  onGenerateShareLink?: () => Promise<string> | string;
  onInvite?: (input: {
    identityMode: InviteIdentityMode;
    displayName?: string;
    email?: string;
    note: string;
  }) => Promise<InviteResult> | InviteResult;
}

export function InviteDialog({
  onClose,
  shareUrl = "",
  onGenerateShareLink,
  onInvite,
}: InviteDialogProps) {
  const [identityMode, setIdentityMode] = useState<InviteIdentityMode>("email");
  const [displayName, setDisplayName] = useState("");
  const [email, setEmail] = useState("");
  const [note, setNote] = useState(
    "I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag."
  );
  const [inviteResult, setInviteResult] = useState<InviteResult | null>(null);
  const [isSending, setIsSending] = useState(false);
  const [isGeneratingShareLink, setIsGeneratingShareLink] = useState(false);
  const [currentShareUrl, setCurrentShareUrl] = useState("");
  const [shareCopied, setShareCopied] = useState(false);
  const [inviteCopied, setInviteCopied] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [shareError, setShareError] = useState<string | null>(null);

  useEffect(() => {
    setCurrentShareUrl(shareUrl);
  }, [shareUrl]);

  useEffect(() => {
    if (!currentShareUrl && onGenerateShareLink) {
      void handleGenerateShareLink();
    }
  }, [currentShareUrl, onGenerateShareLink]);

  const canSubmit = useMemo(() => {
    switch (identityMode) {
      case "email":
        return email.includes("@");
      case "named":
        return displayName.trim().length > 0;
      case "anonymous":
        return true;
      case "preview":
        return true;
      default:
        return false;
    }
  }, [displayName, email, identityMode]);

  const handleGenerateShareLink = async () => {
    if (!onGenerateShareLink || isGeneratingShareLink) {
      return;
    }

    setIsGeneratingShareLink(true);
    setShareError(null);

    try {
      const url = await onGenerateShareLink();
      setCurrentShareUrl(url);
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Failed to generate share link";
      setShareError(message);
    } finally {
      setIsGeneratingShareLink(false);
    }
  };

  const handleSend = async () => {
    if (!onInvite || isSending) {
      return;
    }

    setIsSending(true);
    setError(null);

    try {
      const result = await onInvite({
        identityMode,
        displayName: displayName.trim() || undefined,
        email: email.trim() || undefined,
        note,
      });
      setInviteResult(result);
    } catch (err) {
      const message = err instanceof Error ? err.message : "Failed to create invite";
      setError(message);
    } finally {
      setIsSending(false);
    }
  };

  const handleCopyShareLink = () => {
    if (!currentShareUrl) {
      return;
    }

    navigator.clipboard.writeText(currentShareUrl);
    setShareCopied(true);
    setTimeout(() => setShareCopied(false), 2000);
  };

  const handleCopyInviteLink = () => {
    if (!inviteResult) {
      return;
    }

    navigator.clipboard.writeText(inviteResult.inviteUrl);
    setInviteCopied(true);
    setTimeout(() => setInviteCopied(false), 2000);
  };

  return (
    <div className="dr-invite-overlay" onClick={onClose}>
      <div className="dr-invite-dialog" onClick={(e) => e.stopPropagation()}>
        <div className="dr-invite-dialog__body">
          <p className="dr-invite-dialog__intro">
            Share a reusable article link or mint a tracked reader link for one
            specific reviewer.
          </p>

          <div className="dr-invite-dialog__section">
            <label className="dr-invite-dialog__label">Reusable Article Link</label>
            <div className="dr-invite-dialog__hint">
              Use this when you want one stable review link you can paste into
              email, chat, or docs without creating a tracked reader slot.
            </div>
            <div className="dr-invite-dialog__link-row">
              <input
                className="dr-input dr-invite-dialog__link-input"
                value={currentShareUrl}
                readOnly
                placeholder="Generate a share link"
              />
              <MacButton small onClick={handleGenerateShareLink}>
                {isGeneratingShareLink
                  ? "Generating..."
                  : currentShareUrl
                    ? "Refresh"
                    : "Generate"}
              </MacButton>
              <MacButton
                small
                onClick={handleCopyShareLink}
                disabled={!currentShareUrl}
              >
                {shareCopied ? "Copied!" : "Copy"}
              </MacButton>
            </div>
            {shareError ? (
              <div className="dr-invite-dialog__error">{shareError}</div>
            ) : null}
          </div>

          <div className="dr-invite-dialog__section">
            <label className="dr-invite-dialog__label">Tracked Reader Link</label>
            <div className="dr-invite-dialog__hint">
              Create a unique link for one reader so their review session is
              attributable even when you do not want to use the generic share
              link.
            </div>
            <div className="dr-invite-dialog__mode-row">
              <button
                className={`dr-invite-dialog__mode-btn ${identityMode === "email" ? "dr-invite-dialog__mode-btn--active" : ""}`}
                onClick={() => setIdentityMode("email")}
                type="button"
              >
                Email
              </button>
              <button
                className={`dr-invite-dialog__mode-btn ${identityMode === "named" ? "dr-invite-dialog__mode-btn--active" : ""}`}
                onClick={() => setIdentityMode("named")}
                type="button"
              >
                Named
              </button>
              <button
                className={`dr-invite-dialog__mode-btn ${identityMode === "anonymous" ? "dr-invite-dialog__mode-btn--active" : ""}`}
                onClick={() => setIdentityMode("anonymous")}
                type="button"
              >
                Anonymous
              </button>
            </div>
            {identityMode === "email" && (
              <>
                <label className="dr-invite-dialog__label">Reader Email</label>
                <input
                  className="dr-input"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="reader@example.com"
                />
                <label className="dr-invite-dialog__label">
                  Reader Name
                  <span className="dr-invite-dialog__optional">optional</span>
                </label>
                <input
                  className="dr-input"
                  value={displayName}
                  onChange={(e) => setDisplayName(e.target.value)}
                  placeholder="Workshop Reader"
                />
              </>
            )}
            {identityMode === "named" && (
              <>
                <label className="dr-invite-dialog__label">Reader Name</label>
                <input
                  className="dr-input"
                  value={displayName}
                  onChange={(e) => setDisplayName(e.target.value)}
                  placeholder="Workshop Reader"
                />
                <label className="dr-invite-dialog__label">
                  Reader Email
                  <span className="dr-invite-dialog__optional">optional</span>
                </label>
                <input
                  className="dr-input"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="reader@example.com"
                />
              </>
            )}
            {identityMode === "anonymous" && (
              <div className="dr-invite-dialog__hint dr-invite-dialog__hint--box">
                This creates a one-reader tracked link without a stored name or
                email. The resulting feedback stays attributable to that invite
                slot instead of collapsing into the generic share link.
              </div>
            )}
            <label className="dr-invite-dialog__label">Personal Note</label>
            <textarea
              className="dr-textarea"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              rows={4}
            />
            {inviteResult ? (
              <div className="dr-invite-dialog__sent">
                <div className="dr-invite-dialog__sent-title">
                  Reader Link Created
                </div>
                <div className="dr-invite-dialog__sent-hint">
                  {inviteResult.identityMode === "anonymous"
                    ? "This unique link is associated with an anonymous tracked reader slot."
                    : `This unique link is associated with ${inviteResult.name}${inviteResult.email ? ` (${inviteResult.email})` : ""}.`}
                </div>
                <div className="dr-invite-dialog__link-row">
                  <input
                    className="dr-input dr-invite-dialog__link-input"
                    value={inviteResult.inviteUrl}
                    readOnly
                  />
                  <MacButton small onClick={handleCopyInviteLink}>
                    {inviteCopied ? "Copied!" : "Copy"}
                  </MacButton>
                </div>
              </div>
            ) : null}
            {error ? <div className="dr-invite-dialog__error">{error}</div> : null}
          </div>

          <div className="dr-invite-dialog__actions">
            <MacButton onClick={onClose}>Done</MacButton>
            <MacButton
              primary
              disabled={!canSubmit || isSending}
              onClick={handleSend}
            >
              {isSending
                ? "Creating..."
                : identityMode === "anonymous"
                  ? "Create Anonymous Link"
                  : identityMode === "named"
                    ? "Create Named Link"
                    : "Create Reader Link"}
            </MacButton>
          </div>
        </div>
      </div>
    </div>
  );
}
