import { useState } from "react";
import { REACTION_TYPES } from "../theme/tokens";
import { ReactionBadge } from "../primitives/ReactionBadge";
import { Prose } from "../primitives/Prose";
import type { Reaction, ReactionType } from "../types";
import "./Paragraph.css";

interface ParagraphProps {
  text: string;
  paragraphId: string;
  reactions: Reaction[];
  onReact: (type: ReactionType, text: string) => void;
  onRemoveReaction: (reaction: Reaction) => void;
}

export function Paragraph({
  text,
  paragraphId,
  reactions,
  onReact,
  onRemoveReaction,
}: ParagraphProps) {
  const [activeType, setActiveType] = useState<ReactionType | null>(null);
  const [comment, setComment] = useState("");
  const myReactions = reactions.filter((r) => r.paragraphId === paragraphId);

  const submit = (type: ReactionType) => {
    const reactionText =
      comment.trim() ||
      REACTION_TYPES.find((r) => r.type === type)?.desc ||
      "";
    onReact(type, reactionText);
    setActiveType(null);
    setComment("");
  };

  const handleChipClick = (type: ReactionType) => {
    if (activeType === type) {
      // Submit immediately with default text
      submit(type);
    } else {
      setActiveType(type);
      setComment("");
    }
  };

  return (
    <div
      className={`dr-para ${myReactions.length > 0 ? "dr-para--commented" : ""}`}
    >
      <Prose className="dr-para__text">{text}</Prose>

      {/* Always-visible reaction chips */}
      <div className="dr-para__actions">
        {REACTION_TYPES.map((r) => (
          <span
            key={r.type}
            className={`dr-para__chip ${
              activeType === r.type ? "dr-para__chip--active" : ""
            }`}
            onClick={() => handleChipClick(r.type)}
          >
            {r.icon} {r.label}
          </span>
        ))}
        {myReactions.length > 0 && (
          <span className="dr-para__badge">{myReactions.length}</span>
        )}
      </div>

      {/* Comment input (shown when a chip is selected) */}
      {activeType && (
        <div className="dr-para__comment-box">
          <textarea
            className="dr-para__comment-input"
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                submit(activeType);
              }
              if (e.key === "Escape") {
                setActiveType(null);
                setComment("");
              }
            }}
            placeholder="Add a note (optional) — hit Enter to submit, Esc to cancel"
            rows={2}
            autoFocus
          />
          <div className="dr-para__comment-actions">
            <button
              className="dr-button dr-button--small"
              onClick={() => {
                setActiveType(null);
                setComment("");
              }}
            >
              Cancel
            </button>
            <button
              className="dr-button dr-button--small dr-button--primary"
              onClick={() => submit(activeType)}
            >
              Submit
            </button>
          </div>
        </div>
      )}

      {/* Existing reactions */}
      {myReactions.length > 0 && (
        <div className="dr-para__reactions">
          {myReactions.map((r, i) => (
            <ReactionBadge
              key={i}
              type={r.type}
              text={r.text}
              onRemove={() => onRemoveReaction(r)}
            />
          ))}
        </div>
      )}
    </div>
  );
}
