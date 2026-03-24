import { useState } from "react";
import { TitleBar } from "../chrome/TitleBar";
import { MacButton } from "../chrome/MacButton";
import "./InviteDialog.css";

interface InviteDialogProps {
  onClose: () => void;
  onInvite?: (email: string, note: string) => void;
}

export function InviteDialog({ onClose, onInvite }: InviteDialogProps) {
  const [email, setEmail] = useState("");
  const [note, setNote] = useState(
    "I'm working on an article and would love your honest feedback — especially where you feel confused or where things drag."
  );
  const [sent, setSent] = useState(false);

  const handleSend = () => {
    onInvite?.(email, note);
    setSent(true);
  };

  return (
    <div className="dr-invite-overlay">
      <div className="dr-invite-dialog">
        <TitleBar title="Invite Reader" onClose={onClose} />
        <div className="dr-invite-dialog__body">
          {!sent ? (
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
              <div className="dr-invite-dialog__actions">
                <MacButton onClick={onClose}>Cancel</MacButton>
                <MacButton
                  primary
                  disabled={!email.includes("@")}
                  onClick={handleSend}
                >
                  Send Invite
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
                Your reader will receive a link to review the draft. Remember,
                only about 1 in 4 invitees will engage — that's normal!
              </div>
              <MacButton onClick={onClose}>OK</MacButton>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
