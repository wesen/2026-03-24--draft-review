import { useState, useEffect, useRef, useCallback } from "react";
import type { ReaderArticle, Reaction, ReactionType } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import { WelcomeSplash } from "./WelcomeSplash";
import { ReaderToolbar } from "./ReaderToolbar";
import { SectionView } from "./SectionView";
import { DoneDialog } from "./DoneDialog";
import "./ReaderPage.css";

interface ReaderPageProps {
  article: ReaderArticle;
  /** Optional: pre-existing reactions (e.g. loaded from API) */
  initialReactions?: Reaction[];
  /** Called when a new reaction is added */
  onReactionAdd?: (reaction: Omit<Reaction, "id" | "createdAt">) => void;
  /** Called when a reaction is removed */
  onReactionRemove?: (reaction: Reaction) => void;
}

export function ReaderPage({
  article,
  initialReactions = [],
  onReactionAdd,
  onReactionRemove,
}: ReaderPageProps) {
  const [started, setStarted] = useState(false);
  const [currentSectionId, setCurrentSectionId] = useState(
    article.sections[0]?.id
  );
  const [readSectionIds, setReadSectionIds] = useState<string[]>([]);
  const [reactions, setReactions] = useState<Reaction[]>(initialReactions);
  const [showDone, setShowDone] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  const sectionIndex = article.sections.findIndex(
    (s) => s.id === currentSectionId
  );
  const section = article.sections[sectionIndex];
  const isFirst = sectionIndex === 0;
  const isLast = sectionIndex === article.sections.length - 1;

  const goTo = useCallback((id: string) => {
    setCurrentSectionId(id);
    if (scrollRef.current) scrollRef.current.scrollTop = 0;
  }, []);

  const markRead = useCallback(() => {
    setReadSectionIds((prev) =>
      prev.includes(currentSectionId)
        ? prev
        : [...prev, currentSectionId]
    );
  }, [currentSectionId]);

  const goNext = () => {
    markRead();
    if (!isLast) goTo(article.sections[sectionIndex + 1].id);
  };

  const goPrev = () => {
    if (!isFirst) goTo(article.sections[sectionIndex - 1].id);
  };

  const finish = () => {
    markRead();
    setShowDone(true);
  };

  const handleReact = (
    paragraphId: string,
    type: ReactionType,
    text: string
  ) => {
    const newReaction: Reaction = {
      id: `local-${Date.now()}`,
      articleId: article.id,
      sectionId: currentSectionId,
      paragraphId,
      readerId: "anonymous",
      readerName: "You",
      type,
      text,
      createdAt: new Date().toISOString(),
    };
    setReactions((prev) => [...prev, newReaction]);
    onReactionAdd?.({
      articleId: article.id,
      sectionId: currentSectionId,
      paragraphId,
      readerId: "anonymous",
      readerName: "You",
      type,
      text,
    });
  };

  const handleRemoveReaction = (reaction: Reaction) => {
    setReactions((prev) => prev.filter((r) => r !== reaction));
    onReactionRemove?.(reaction);
  };

  // Keyboard navigation
  useEffect(() => {
    if (!started) return;
    const handler = (e: KeyboardEvent) => {
      if (e.key === "ArrowRight" || e.key === "ArrowDown") {
        e.preventDefault();
        if (!isLast) goNext();
      }
      if (e.key === "ArrowLeft" || e.key === "ArrowUp") {
        e.preventDefault();
        if (!isFirst) goPrev();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  });

  // Stats for done dialog
  const stats: Record<string, number> = {};
  REACTION_TYPES.forEach((r) => {
    stats[r.type] = reactions.filter((rx) => rx.type === r.type).length;
  });

  if (!started) {
    return (
      <WelcomeSplash
        article={article}
        onStart={() => {
          setStarted(true);
          setCurrentSectionId(article.sections[0]?.id);
        }}
      />
    );
  }

  return (
    <div className="dr-reader-page">
      <ReaderToolbar
        sections={article.sections}
        currentSectionId={currentSectionId}
        readSectionIds={readSectionIds}
        onPickSection={goTo}
      />

      <div className="dr-reader-page__content" ref={scrollRef}>
        {section && (
          <>
            <SectionView
              section={section}
              sectionIndex={sectionIndex}
              totalSections={article.sections.length}
              reactions={reactions}
              onReact={handleReact}
              onRemoveReaction={handleRemoveReaction}
            />

            <div className="dr-reader-page__nav">
              <MacButton onClick={goPrev} disabled={isFirst}>
                ← Previous
              </MacButton>
              {isLast ? (
                <MacButton primary onClick={finish}>
                  Finish Review ✓
                </MacButton>
              ) : (
                <MacButton primary onClick={goNext}>
                  Next Section →
                </MacButton>
              )}
            </div>

            <div className="dr-reader-page__hint">
              Use arrow keys to navigate · hover paragraphs to react
            </div>
          </>
        )}
      </div>

      {showDone && (
        <DoneDialog stats={stats} onClose={() => setShowDone(false)} />
      )}
    </div>
  );
}
