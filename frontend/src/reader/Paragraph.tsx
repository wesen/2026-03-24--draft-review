import { useState } from "react";
import { ReactionPicker } from "../primitives/ReactionPicker";
import { ReactionBadge } from "../primitives/ReactionBadge";
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
  const [hovered, setHovered] = useState(false);
  const [picking, setPicking] = useState(false);
  const myReactions = reactions.filter((r) => r.paragraphId === paragraphId);

  return (
    <div
      className="dr-paragraph"
      onMouseEnter={() => setHovered(true)}
      onMouseLeave={() => setHovered(false)}
    >
      <p
        className={`dr-paragraph__text ${
          myReactions.length > 0 ? "dr-paragraph__text--reacted" : ""
        }`}
      >
        {text}
      </p>

      {hovered && !picking && (
        <div className="dr-paragraph__add" onClick={() => setPicking(true)}>
          +
        </div>
      )}

      {picking && (
        <div className="dr-paragraph__picker">
          <ReactionPicker
            onSubmit={(type, commentText) => {
              onReact(type, commentText);
              setPicking(false);
            }}
            onCancel={() => setPicking(false)}
          />
        </div>
      )}

      {myReactions.length > 0 && (
        <div className="dr-paragraph__reactions">
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
