import { useState, useEffect, useRef } from "react";
import { REACTION_TYPES } from "../theme/tokens";
import type { ReactionType } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import "./ReactionPicker.css";

interface ReactionPickerProps {
  onSubmit: (type: ReactionType, text: string) => void;
  onCancel: () => void;
}

export function ReactionPicker({ onSubmit, onCancel }: ReactionPickerProps) {
  const [chosenType, setChosenType] = useState<ReactionType | null>(null);
  const [comment, setComment] = useState("");
  const inputRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (chosenType && inputRef.current) inputRef.current.focus();
  }, [chosenType]);

  const submit = () => {
    if (!chosenType) return;
    const text =
      comment.trim() ||
      REACTION_TYPES.find((r) => r.type === chosenType)?.desc ||
      "";
    onSubmit(chosenType, text);
  };

  return (
    <div className="dr-reaction-picker">
      <div className="dr-reaction-picker__types">
        {REACTION_TYPES.map((r) => (
          <div
            key={r.type}
            onClick={() => setChosenType(r.type)}
            className={`dr-reaction-picker__type ${
              chosenType === r.type ? "dr-reaction-picker__type--selected" : ""
            }`}
          >
            <span className="dr-reaction-picker__type-icon">{r.icon}</span>
            {r.label}
          </div>
        ))}
      </div>
      {chosenType && (
        <>
          <textarea
            ref={inputRef}
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                submit();
              }
            }}
            placeholder="Add a note (optional) — what specifically?"
            rows={2}
            className="dr-reaction-picker__input"
          />
          <div className="dr-reaction-picker__actions">
            <MacButton small onClick={onCancel}>
              Cancel
            </MacButton>
            <MacButton small primary onClick={submit}>
              Submit
            </MacButton>
          </div>
        </>
      )}
    </div>
  );
}
