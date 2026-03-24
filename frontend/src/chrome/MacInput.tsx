import type { InputHTMLAttributes } from "react";
import "./MacInput.css";

interface MacInputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
}

export function MacInput({ label, id, className = "", ...props }: MacInputProps) {
  return (
    <div className="dr-field">
      {label && (
        <label className="dr-field__label" htmlFor={id}>
          {label}
        </label>
      )}
      <input className={`dr-input ${className}`} id={id} {...props} />
    </div>
  );
}
