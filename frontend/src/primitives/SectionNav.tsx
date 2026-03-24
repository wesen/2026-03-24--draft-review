import type { Section } from "../types";
import "./SectionNav.css";

interface SectionNavProps {
  sections: Section[];
  currentId: string;
  readIds: string[];
  onPick: (id: string) => void;
}

export function SectionNav({ sections, currentId, readIds, onPick }: SectionNavProps) {
  return (
    <div className="dr-section-nav">
      {sections.map((s) => {
        const isCurrent = s.id === currentId;
        const isRead = readIds.includes(s.id);
        return (
          <div
            key={s.id}
            title={s.title}
            onClick={() => onPick(s.id)}
            className={[
              "dr-section-nav__dot",
              isCurrent && "dr-section-nav__dot--current",
              isRead && "dr-section-nav__dot--read",
            ]
              .filter(Boolean)
              .join(" ")}
          />
        );
      })}
    </div>
  );
}
