import "./ProgressBar.css";

interface ProgressBarProps {
  percent: number;
  width?: string | number;
  height?: number;
}

export function ProgressBar({ percent, width = "100%", height = 12 }: ProgressBarProps) {
  return (
    <div className="dr-progress" style={{ width, height }}>
      <div
        className="dr-progress__fill"
        style={{ width: `${Math.min(100, Math.max(0, percent))}%` }}
      />
    </div>
  );
}
