import { useState, useEffect, useRef, useCallback } from "react";
import type { ReaderArticle, Reaction, ReactionType } from "../types";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import {
  useAddReviewReactionMutation,
  useStartReviewMutation,
  useSubmitReviewSummaryMutation,
  useUpdateReviewProgressMutation,
} from "../api/readerApi";
import { useAppDispatch, useAppSelector } from "../store/hooks";
import {
  startReading,
  setSessionId as setSessionIdAction,
  goToSection,
  markSectionRead,
  resetReader,
} from "../store/readerSlice";
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
  reviewToken?: string;
  /** When true, reactions are disabled */
  readOnly?: boolean;
  /** When true, show preview-specific affordances */
  previewMode?: boolean;
  /** Called when user clicks "Back to Editor" in preview mode */
  onBackToEditor?: () => void;
}

export function ReaderPage({
  article,
  initialReactions = [],
  onReactionAdd,
  onReactionRemove,
  reviewToken,
  readOnly = false,
  previewMode = false,
  onBackToEditor,
}: ReaderPageProps) {
  const dispatch = useAppDispatch();

  // Redux state
  const started = useAppSelector((s) => s.reader.started);
  const sessionId = useAppSelector((s) => s.reader.sessionId);
  const currentSectionId =
    useAppSelector((s) => s.reader.currentSectionId) ?? article.sections[0]?.id;
  const readSectionIds = useAppSelector((s) => s.reader.readSectionIds);

  // Local state (optimistic reactions + dialog toggle)
  const [reactions, setReactions] = useState<Reaction[]>(initialReactions);
  const [showDone, setShowDone] = useState(false);

  const scrollRef = useRef<HTMLDivElement>(null);
  const [startReview] = useStartReviewMutation();
  const [updateReviewProgress] = useUpdateReviewProgressMutation();
  const [addReviewReaction] = useAddReviewReactionMutation();
  const [submitReviewSummary] = useSubmitReviewSummaryMutation();

  // Reset reader state on unmount
  useEffect(() => {
    return () => {
      dispatch(resetReader());
    };
  }, [dispatch]);

  const sectionIndex = article.sections.findIndex(
    (s) => s.id === currentSectionId
  );
  const section = article.sections[sectionIndex];
  const isFirst = sectionIndex === 0;
  const isLast = sectionIndex === article.sections.length - 1;

  const goTo = useCallback(
    (id: string) => {
      dispatch(goToSection(id));
      if (scrollRef.current) scrollRef.current.scrollTop = 0;
    },
    [dispatch]
  );

  const persistProgress = useCallback(
    async (sectionId: string) => {
      if (!sessionId) return;

      const nextReadCount = readSectionIds.includes(sectionId)
        ? readSectionIds.length
        : readSectionIds.length + 1;
      const progressPercent = Math.round(
        (nextReadCount / article.sections.length) * 100
      );

      try {
        await updateReviewProgress({
          sessionId,
          sectionId,
          paragraphId: `${sectionId}-p0`,
          completed: true,
          progressPercent,
        }).unwrap();
      } catch (error) {
        console.error("Failed to persist review progress", error);
      }
    },
    [article.sections.length, readSectionIds, sessionId, updateReviewProgress]
  );

  const markRead = useCallback(
    (sectionId: string) => {
      dispatch(markSectionRead(sectionId));
      void persistProgress(sectionId);
    },
    [dispatch, persistProgress]
  );

  const goNext = useCallback(() => {
    markRead(currentSectionId);
    if (!isLast) goTo(article.sections[sectionIndex + 1].id);
  }, [currentSectionId, isLast, sectionIndex, article.sections, markRead, goTo]);

  const goPrev = useCallback(() => {
    if (!isFirst) goTo(article.sections[sectionIndex - 1].id);
  }, [isFirst, sectionIndex, article.sections, goTo]);

  const finish = async () => {
    markRead(currentSectionId);
    if (sessionId) {
      try {
        await submitReviewSummary({
          sessionId,
          notifyNewVersion: false,
        }).unwrap();
      } catch (error) {
        console.error("Failed to submit review summary", error);
      }
    }
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

    if (sessionId) {
      void addReviewReaction({
        sessionId,
        sectionId: currentSectionId,
        paragraphId,
        type,
        text,
      }).unwrap().catch((error) => {
        console.error("Failed to persist review reaction", error);
      });
    }
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
  }, [started, isFirst, isLast, goNext, goPrev]);

  // Stats for done dialog
  const stats: Record<string, number> = {};
  REACTION_TYPES.forEach((r) => {
    stats[r.type] = reactions.filter((rx) => rx.type === r.type).length;
  });

  if (!started) {
    return (
      <WelcomeSplash
        article={article}
        onStart={async () => {
          if (reviewToken && !sessionId) {
            try {
              const result = await startReview({ token: reviewToken }).unwrap();
              dispatch(setSessionIdAction(result.session.id));
            } catch (error) {
              console.error("Failed to start review session", error);
            }
          }
          dispatch(startReading(article.sections[0]?.id));
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

      {previewMode && (
        <div className="dr-reader-page__preview-banner">
          <span>
            {readOnly
              ? "Preview mode \u2014 reactions are disabled"
              : "Preview mode \u2014 test reactions stay local to this preview session"}
          </span>
          {onBackToEditor && (
            <MacButton small onClick={onBackToEditor}>
              {"\u2190"} Back to Editor
            </MacButton>
          )}
        </div>
      )}

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
              readOnly={readOnly}
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
