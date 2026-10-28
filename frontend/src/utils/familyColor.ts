import type { Family } from "../types";

export interface FamilyColor {
  /** Pale tint background — DESIGN.md's card-tint-* palette. */
  bg: string;
  /** Saturated accent — left border / swatch / icon. */
  border: string;
  /** Deeper same-hue tone — readable as text on `bg`. */
  text: string;
}

// Straight from DESIGN.md's card-tint-* + brand-* tokens — the same
// pastel-tint/accent pairing the spec itself uses for badge-tag-*
// components. Six is enough to stay visually distinct; beyond that it
// cycles, same as any categorical color scale.
const PALETTE: FamilyColor[] = [
  { bg: "#ffe8d4", border: "#dd5b00", text: "#793400" }, // peach / orange
  { bg: "#fde0ec", border: "#ff64c8", text: "#a02e6d" }, // rose / pink
  { bg: "#d9f3e1", border: "#1aae39", text: "#1aae39" }, // mint / green
  { bg: "#e6e0f5", border: "#7b3ff2", text: "#391c57" }, // lavender / purple
  { bg: "#dcecfa", border: "#0075de", text: "#005bab" }, // sky / blue
  { bg: "#fef7d6", border: "#523410", text: "#523410" }, // yellow / brown
];

/**
 * Assigns each family a stable color for the session, by its position in
 * the user's family list (not a hash of the id) — simplest thing that
 * keeps a given family the same color across every page without needing
 * to persist anything.
 */
export function familyColorMap(families: Family[]): Map<string, FamilyColor> {
  const map = new Map<string, FamilyColor>();
  families.forEach((f, i) => map.set(f.id, PALETTE[i % PALETTE.length]));
  return map;
}
