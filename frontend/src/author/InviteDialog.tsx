import { useState } from "react";
import { TitleBar } from "../chrome/TitleBar";
import { MacButton } from "../chrome/MacButton";
import "./InviteDialog.css";

interface InviteResult {
  email: string;
  inviteUrl: string;
}

interface InviteDialogProps {
  onClose: () => void;
  shareUrl?: string;
  onGenerateShareLink?: () => Promise<string> | string;
  onInvite?: (email: string, note: string) => Promise<InviteResult> | InviteResult;
}

export function InviteDialog({
  onClose,
  shareUrl,
  onGenerateShareLink,
  onInvite,
}: InviteDialogProps) {
  const [email, setEmail] = useState("");
  const [note, setNote] = useState(
    "I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag."
  );
  const [inviteResult, setInviteResult] = useState<InviteResult | null>(null);
  const [isSending, setIsSending] = useState(false);
  const [isGeneratingShareLink, setIsGeneratingShareLink] = useState(false);
  const [currentShareUrl, setCurrentShareUrl] = useState(shareUrl ?? "");
  const [shareCopied, setShareCopied] = useState(false);
  const [inviteCopied, setInviteCopied] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [shareError, setShareError] = useState<string | null>(null);

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
      const result = await onInvite(email, note);
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
    <div className="dr-invite-overlay">
      <div className="dr-invite-dialog">
        <TitleBar title="Invite Reader" onClose={onClose} />
        <div className="dr-invite-dialog__body">
          <p className="dr-invite-dialog__intro">
            Share a generic article link or create a reader-specific link tied
            to an email address.
          </p>

          <div className="dr-invite-dialog__section">
            <label className="dr-invite-dialog__label">Share Link</label>
            <div className="dr-invite-dialog__hint">
              Use this when you just want a reusable review link without tying
              it to a specific reader email.
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
            <label className="dr-invite-dialog__label">
              Reader Email
              <span className="dr-invite-dialog__optional">optional</span>
            </label>
            <div className="dr-invite-dialog__hint">
              Create a distinct reader link so you can track who it was meant
              for.
            </div>
            <input
              className="dr-input"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="reader@example.com"
            />
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
                  This unique link is associated with {inviteResult.email}.
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
              disabled={!email.includes("@") || isSending}
              onClick={handleSend}
            >
              {isSending ? "Creating..." : "Create Reader Link"}
            </MacButton>
          </div>
        </div>
      </div>
    </div>
  );
}
