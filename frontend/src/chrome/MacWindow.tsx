import type { ReactNode } from "react";
import "./MacWindow.css";

interface MacWindowProps {
  title?: string;
  children: ReactNode;
  maximized?: boolean;
  onClose?: () => void;
  zIndex?: number;
}

export function MacWindow({
  children,
  maximized = false,
  zIndex = 1,
}: MacWindowProps) {
  const className = maximized ? "dr-window dr-window--maximized" : "dr-window";

  return (
    <div className={className} style={{ zIndex }}>
      <div className="dr-window__content">{children}</div>
    </div>
  );
}
