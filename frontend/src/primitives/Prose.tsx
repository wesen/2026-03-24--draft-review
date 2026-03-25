import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import "./Prose.css";

interface ProseProps {
  children: string;
  className?: string;
}

/**
 * Renders a markdown string as styled prose.
 * Used in the reader view for paragraph content.
 */
export function Prose({ children, className = "" }: ProseProps) {
  return (
    <div className={`dr-prose ${className}`}>
      <ReactMarkdown remarkPlugins={[remarkGfm]}>{children}</ReactMarkdown>
    </div>
  );
}
