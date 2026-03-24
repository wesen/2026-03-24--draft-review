/**
 * Design token constants that mirror the CSS custom properties.
 * Use these in TypeScript when you need token values programmatically
 * (e.g., for canvas drawing or computed styles).
 * For component styling, prefer the CSS custom properties in mac-os-1.css.
 */

export const FONTS = {
  title: `"Chicago_12", "ChicagoFLF", "Geneva", monospace`,
  body: `"Geneva", "Monaco", monospace`,
} as const;

export const REACTION_TYPES = [
  { type: "useful", icon: "\u2605", label: "Useful!", desc: "I can use this" },
  { type: "confused", icon: "?", label: "Confused", desc: "I don't follow" },
  { type: "slow", icon: "\u25CE", label: "Slow", desc: "This drags" },
  { type: "favorite", icon: "\u2665", label: "Love it", desc: "Highlight" },
] as const;

export type ReactionType = (typeof REACTION_TYPES)[number]["type"];
