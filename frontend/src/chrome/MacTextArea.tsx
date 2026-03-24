import type { TextareaHTMLAttributes } from "react";
import "./MacInput.css";

interface MacTextAreaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
}

export function MacTextArea({ label, id, className = "", ...props }: MacTextAreaProps) {
  return (
    <div className="dr-field">
      {label && (
        <label className="dr-field__label" htmlFor={id}>
          {label}
        </label>
      )}
      <textarea className={`dr-textarea ${className}`} id={id} {...props} />
    </div>
  );
}
