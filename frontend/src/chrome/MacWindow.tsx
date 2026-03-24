import { useState, useEffect, useRef, type ReactNode } from "react";
import { TitleBar } from "./TitleBar";
import "./MacWindow.css";

interface MacWindowProps {
  title: string;
  children: ReactNode;
  x?: number;
  y?: number;
  w?: number;
  h?: number;
  maximized?: boolean;
  onClose?: () => void;
  zIndex?: number;
}

export function MacWindow({
  title,
  children,
  x = 40,
  y = 40,
  w = 600,
  h = 450,
  maximized = false,
  onClose,
  zIndex = 1,
}: MacWindowProps) {
  const [pos, setPos] = useState({ x, y });
  const [dragging, setDragging] = useState(false);
  const dragRef = useRef<{ startX: number; startY: number } | null>(null);

  const onMouseDown = (e: React.MouseEvent) => {
    if (maximized) return;
    dragRef.current = { startX: e.clientX - pos.x, startY: e.clientY - pos.y };
    setDragging(true);
  };

  useEffect(() => {
    if (!dragging) return;
    const onMove = (e: MouseEvent) => {
      if (dragRef.current) {
        setPos({
          x: e.clientX - dragRef.current.startX,
          y: e.clientY - dragRef.current.startY,
        });
      }
    };
    const onUp = () => setDragging(false);
    window.addEventListener("mousemove", onMove);
    window.addEventListener("mouseup", onUp);
    return () => {
      window.removeEventListener("mousemove", onMove);
      window.removeEventListener("mouseup", onUp);
    };
  }, [dragging]);

  const className = maximized ? "dr-window dr-window--maximized" : "dr-window";
  const style = maximized
    ? { zIndex }
    : { left: pos.x, top: pos.y, width: w, height: h, zIndex };

  return (
    <div className={className} style={style}>
      <div
        onMouseDown={onMouseDown}
        style={{ cursor: maximized ? "default" : "grab" }}
      >
        <TitleBar title={title} onClose={onClose} />
      </div>
      <div className="dr-window__content">{children}</div>
    </div>
  );
}
