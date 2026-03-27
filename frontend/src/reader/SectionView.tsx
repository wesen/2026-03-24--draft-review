import type { Section, Reaction, ReactionType } from "../types";
import { deriveSectionBlocks } from "../lib/markdownBlocks";
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
  readOnly?: boolean;
}

export function SectionView({
  section,
  sectionIndex,
  totalSections,
  reactions,
  onReact,
  onRemoveReaction,
  readOnly = false,
}: SectionViewProps) {
  const sectionReactions = reactions.filter(
    (r) => r.sectionId === section.id
  );
  const blocks = deriveSectionBlocks(section);

  return (
    <div className="dr-section-view">
      <div className="dr-section-view__divider">
        {"\u00A7"} {section.title.toUpperCase()}
      </div>
      <div className="dr-section-view__header">
        <span className="dr-section-view__badge">
          {sectionIndex + 1}/{totalSections}
        </span>
        <h2 className="dr-section-view__title">{section.title}</h2>
      </div>

      {blocks.map((block) => (
        <Paragraph
          key={block.id}
          text={block.markdown}
          paragraphId={block.id}
          reactions={sectionReactions}
          onReact={(type, reactionText) =>
            onReact(block.id, type, reactionText)
          }
          onRemoveReaction={onRemoveReaction}
          readOnly={readOnly}
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
