import { REACTION_TYPES } from "../theme/tokens";
import { TitleBar } from "../chrome/TitleBar";
import { MacButton } from "../chrome/MacButton";
import "./DoneDialog.css";

interface DoneDialogProps {
  stats: Record<string, number>;
  onClose: () => void;
}

export function DoneDialog({ stats, onClose }: DoneDialogProps) {
  return (
    <div className="dr-done-overlay">
      <div className="dr-done-dialog">
        <TitleBar title="Review Complete" onClose={onClose} />
        <div className="dr-done-dialog__body">
          <div className="dr-done-dialog__emoji">🎉</div>
          <div className="dr-done-dialog__heading">Thank you for reading!</div>
          <div className="dr-done-dialog__grid">
            {REACTION_TYPES.map((r) => (
              <div key={r.type} className="dr-done-dialog__stat">
                <div className="dr-done-dialog__stat-icon">{r.icon}</div>
                <div className="dr-done-dialog__stat-count">
                  {stats[r.type] || 0}
                </div>
                <div className="dr-done-dialog__stat-label">{r.label}</div>
              </div>
            ))}
          </div>
          <p className="dr-done-dialog__message">
            Your reactions have been saved. The author will use them to improve
            the next revision. Every piece of feedback makes the final article
            better.
          </p>
          <MacButton primary onClick={onClose}>
            Done
          </MacButton>
        </div>
      </div>
    </div>
  );
}
