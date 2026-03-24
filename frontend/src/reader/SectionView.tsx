import type { Section, Reaction, ReactionType } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { Paragraph } from "./Paragraph";
import "./SectionView.css";

interface SectionViewProps {
  section: Section;
  sectionIndex: number;
  totalSections: number;
  reactions: Reaction[];
  onReact: (paragraphId: string, type: ReactionType, text: string) => void;
  onRemoveReaction: (reaction: Reaction) => void;
}

export function SectionView({
  section,
  sectionIndex,
  totalSections,
  reactions,
  onReact,
  onRemoveReaction,
}: SectionViewProps) {
  const sectionReactions = reactions.filter(
    (r) => r.sectionId === section.id
  );

  return (
    <div className="dr-section-view">
      <div className="dr-section-view__header">
        <span className="dr-section-view__badge">
          {sectionIndex + 1} / {totalSections}
        </span>
        <h2 className="dr-section-view__title">{section.title}</h2>
      </div>
      <div className="dr-section-view__divider" />

      {section.paragraphs.map((text, pi) => (
        <Paragraph
          key={`${section.id}-${pi}`}
          text={text}
          paragraphId={`${section.id}-p${pi}`}
          reactions={sectionReactions}
          onReact={(type, reactionText) =>
            onReact(`${section.id}-p${pi}`, type, reactionText)
          }
          onRemoveReaction={onRemoveReaction}
        />
      ))}

      {sectionReactions.length > 0 && (
        <div className="dr-section-view__summary">
          <span className="dr-section-view__summary-label">This section:</span>
          {REACTION_TYPES.map((r) => {
            const count = sectionReactions.filter(
              (rx) => rx.type === r.type
            ).length;
            return count > 0 ? (
              <span key={r.type}>
                {r.icon} {count}
              </span>
            ) : null;
          })}
        </div>
      )}
    </div>
  );
}
