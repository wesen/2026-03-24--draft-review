export interface Section {
  id: string;
  title: string;
  bodyMarkdown: string;
}

export interface Article {
  id: string;
  title: string;
  author: string;
  version: string;
  status: "draft" | "in_review" | "complete" | "archived";
  intro: string;
  shareUrl?: string;
  sections: Section[];
  createdAt: string;
  updatedAt: string;
}

export interface ReaderArticle {
  id: string;
  title: string;
  author: string;
  version: string;
  intro: string;
  sections: Section[];
}
