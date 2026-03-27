import type { Section } from "../types";

export interface SectionBlock {
  id: string;
  index: number;
  markdown: string;
  plainText: string;
}

export function normalizeMarkdownBody(bodyMarkdown: string): string {
  return bodyMarkdown.replace(/\r\n/g, "\n").replace(/\r/g, "\n");
}

export function splitMarkdownIntoBlocks(bodyMarkdown: string): string[] {
  const normalized = normalizeMarkdownBody(bodyMarkdown).replace(/^\n+|\n+$/g, "");
  if (normalized.trim() === "") {
    return [""];
  }

  return normalized
    .split(/\n{2,}/)
    .map((block) => block.replace(/^\n+|\n+$/g, ""))
    .filter((block) => block.length > 0);
}

export function markdownToPlainText(markdown: string): string {
  return markdown
    .replace(/!\[([^\]]*)\]\([^)]+\)/g, "$1")
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1")
    .replace(/`{1,3}([^`]+)`{1,3}/g, "$1")
    .replace(/^>\s?/gm, "")
    .replace(/^#{1,6}\s+/gm, "")
    .replace(/[*_~]/g, "")
    .replace(/\n+/g, " ")
    .replace(/\s+/g, " ")
    .trim();
}

export function deriveSectionBlocks(section: Section): SectionBlock[] {
  return splitMarkdownIntoBlocks(section.bodyMarkdown).map((markdown, index) => ({
    id: `${section.id}-p${index}`,
    index,
    markdown,
    plainText: markdownToPlainText(markdown),
  }));
}

export function getBlockIndexFromParagraphId(paragraphId: string): number | null {
  const match = paragraphId.match(/-p(\d+)$/);
  if (!match) {
    return null;
  }

  const index = Number.parseInt(match[1], 10);
  return Number.isNaN(index) ? null : index;
}

export function estimateMarkdownWordCount(markdown: string): number {
  const plainText = markdownToPlainText(markdown);
  if (plainText === "") {
    return 0;
  }

  return plainText.split(/\s+/).filter(Boolean).length;
}
