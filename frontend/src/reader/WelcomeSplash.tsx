import type { ReaderArticle } from "../types";
import { estimateMarkdownWordCount } from "../lib/markdownBlocks";
import { REACTION_TYPES } from "../theme/tokens";
import { MacButton } from "../chrome/MacButton";
import "./WelcomeSplash.css";

interface WelcomeSplashProps {
  article: ReaderArticle;
  onStart: () => void;
}

export function WelcomeSplash({ article, onStart }: WelcomeSplashProps) {
  const wordCount = article.sections.reduce(
    (a, s) => a + estimateMarkdownWordCount(s.bodyMarkdown),
    0
  );
  const readTime = Math.ceil(wordCount / 200);

  return (
    <div className="dr-welcome">
      <div className="dr-welcome__inner">
        <div className="dr-welcome__icon">{"\u270E"}</div>
        <div className="dr-welcome__title">{article.title}</div>
        <div className="dr-welcome__author">by {article.author}</div>
        <div className="dr-welcome__version">{article.version}</div>

        <div className="dr-welcome__note">
          <div className="dr-welcome__note-header">A note from the author:</div>
          {article.intro}
        </div>

        <div className="dr-welcome__howto">
          <div className="dr-welcome__howto-header">How to leave feedback:</div>
          <div className="dr-welcome__howto-grid">
            {REACTION_TYPES.map((r) => (
              <div key={r.type} className="dr-welcome__howto-row">
                <span className="dr-welcome__howto-icon">{r.icon}</span>
                <span>
                  <strong>{r.label}</strong> — {r.desc}
                </span>
              </div>
            ))}
          </div>
          <div className="dr-welcome__howto-hint">
            Hover any paragraph and click a <strong>reaction chip</strong> to
            react. A short note explaining <em>why</em> is even more helpful.
            Don't worry about typos or grammar — focus on what's useful,
            confusing, or slow.
          </div>
        </div>

        <MacButton primary onClick={onStart}>
          Begin Reading →
        </MacButton>
        <div className="dr-welcome__meta">
          {article.sections.length} sections · ~{wordCount} words · {readTime}{" "}
          min read
        </div>
      </div>
    </div>
  );
}
