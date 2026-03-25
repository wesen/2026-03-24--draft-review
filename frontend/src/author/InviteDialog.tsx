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
  onInvite?: (email: string, note: string) => Promise<InviteResult> | InviteResult;
}

export function InviteDialog({ onClose, onInvite }: InviteDialogProps) {
  const [email, setEmail] = useState("");
  const [note, setNote] = useState(
    "I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag."
  );
  const [inviteResult, setInviteResult] = useState<InviteResult | null>(null);
  const [isSending, setIsSending] = useState(false);
  const [copied, setCopied] = useState(false);
  const [error, setError] = useState<string | null>(null);

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

  const handleCopy = () => {
    if (!inviteResult) {
      return;
    }

    navigator.clipboard.writeText(inviteResult.inviteUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="dr-invite-overlay">
      <div className="dr-invite-dialog">
        <TitleBar title="Invite Reader" onClose={onClose} />
        <div className="dr-invite-dialog__body">
          {!inviteResult ? (
            <>
              <p className="dr-invite-dialog__intro">
                Send a personal invitation. Remember: you want readers who would
                genuinely read this topic, even if you didn't write it.
              </p>
              <label className="dr-invite-dialog__label">Email Address</label>
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
              {error ? <div className="dr-invite-dialog__error">{error}</div> : null}
              <div className="dr-invite-dialog__actions">
                <MacButton onClick={onClose}>Cancel</MacButton>
                <MacButton
                  primary
                  disabled={!email.includes("@") || isSending}
                  onClick={handleSend}
                >
                  {isSending ? "Sending..." : "Send Invite"}
                </MacButton>
              </div>
            </>
          ) : (
            <div className="dr-invite-dialog__sent">
              <div style={{ fontSize: 32, marginBottom: 8 }}>{"\u2709"}</div>
              <div className="dr-invite-dialog__sent-title">
                Invitation Sent!
              </div>
              <div className="dr-invite-dialog__sent-hint">
                This unique link is associated with {inviteResult.email}. You can
                paste it into email, Slack, or anywhere else you normally reach
                that reader.
              </div>
              <label className="dr-invite-dialog__label">Reader Link</label>
              <div className="dr-invite-dialog__link-row">
                <input
                  className="dr-input dr-invite-dialog__link-input"
                  value={inviteResult.inviteUrl}
                  readOnly
                />
                <MacButton small onClick={handleCopy}>
                  {copied ? "Copied!" : "Copy Link"}
                </MacButton>
              </div>
              <div className="dr-invite-dialog__actions dr-invite-dialog__actions--sent">
                <MacButton onClick={handleCopy}>
                  {copied ? "Copied!" : "Copy Link"}
                </MacButton>
                <MacButton primary onClick={onClose}>Done</MacButton>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
