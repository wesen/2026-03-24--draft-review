import type { ReactionType } from "../theme/tokens";

export type { ReactionType };

export interface Reaction {
  id: string;
  articleId: string;
  sectionId: string;
  paragraphId: string;
  readerId: string;
  readerName: string;
  type: ReactionType;
  text: string;
  createdAt: string;
}

export interface AddReactionDto {
  articleId: string;
  sectionId: string;
  paragraphId: string;
  type: ReactionType;
  text: string;
  readerToken?: string;
}
