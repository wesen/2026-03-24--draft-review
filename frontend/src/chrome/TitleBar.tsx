import type { ReactNode } from "react";

interface TitleBarProps {
  title: string;
  onClose?: () => void;
  extra?: ReactNode;
}

export function TitleBar({ title, onClose, extra }: TitleBarProps) {
  return (
    <div className="dr-titlebar">
      {onClose && (
        <div className="dr-titlebar__close" onClick={onClose} />
      )}
      <div className="dr-titlebar__stripe" />
      <span className="dr-titlebar__title">{title}</span>
      <div className="dr-titlebar__stripe" />
      {extra}
    </div>
  );
}
