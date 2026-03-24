import type { ReactNode } from "react";
import "./MacButton.css";

interface MacButtonProps {
  children: ReactNode;
  onClick?: () => void;
  primary?: boolean;
  small?: boolean;
  disabled?: boolean;
  fullWidth?: boolean;
  className?: string;
}

export function MacButton({
  children,
  onClick,
  primary = false,
  small = false,
  disabled = false,
  fullWidth = false,
  className = "",
}: MacButtonProps) {
  const classes = [
    "dr-button",
    primary && "dr-button--primary",
    small && "dr-button--small",
    fullWidth && "dr-button--full",
    className,
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <button className={classes} disabled={disabled} onClick={onClick}>
      {children}
    </button>
  );
}
