import { REACTION_TYPES } from "../theme/tokens";
import type { ReactionType } from "../theme/tokens";
import "./ReactionBadge.css";

interface ReactionBadgeProps {
  type: ReactionType;
  text: string;
  onRemove?: () => void;
}

export function ReactionBadge({ type, text, onRemove }: ReactionBadgeProps) {
  const reaction = REACTION_TYPES.find((r) => r.type === type);
  return (
    <div className="dr-reaction-badge">
      <span className="dr-reaction-badge__icon">{reaction?.icon}</span>
      <span className="dr-reaction-badge__text">{text}</span>
      {onRemove && (
        <span className="dr-reaction-badge__remove" onClick={onRemove}>
          {"\u2715"}
        </span>
      )}
    </div>
  );
}
