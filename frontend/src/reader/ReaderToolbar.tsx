import type { Section } from "../types";
import { SectionNav } from "../primitives/SectionNav";
import { ProgressBar } from "../primitives/ProgressBar";
import "./ReaderToolbar.css";

interface ReaderToolbarProps {
  sections: Section[];
  currentSectionId: string;
  readSectionIds: string[];
  onPickSection: (id: string) => void;
}

export function ReaderToolbar({
  sections,
  currentSectionId,
  readSectionIds,
  onPickSection,
}: ReaderToolbarProps) {
  const progress = (readSectionIds.length / sections.length) * 100;

  return (
    <div className="dr-reader-toolbar">
      <SectionNav
        sections={sections}
        currentId={currentSectionId}
        readIds={readSectionIds}
        onPick={onPickSection}
      />
      <div className="dr-reader-toolbar__bar">
        <ProgressBar percent={progress} height={8} />
      </div>
      <span className="dr-reader-toolbar__percent">
        {Math.round(progress)}%
      </span>
    </div>
  );
}
