import { baseApi } from "./baseApi";
import type { ReaderArticle, Reaction, ReactionType } from "../types";

interface ReaderLinkResponse {
  reader: { id: string; name: string };
  article: ReaderArticle;
}

interface ReviewSession {
  id: string;
  articleId: string;
  articleVersionId: string;
  readerId: string;
  readerName: string;
  readerEmail?: string;
  isAnonymous: boolean;
  progressPercent: number;
  startedAt: string;
  lastActiveAt: string;
  completedAt?: string;
}

interface StartReviewResponse extends ReaderLinkResponse {
  session: ReviewSession;
}

interface ReviewProgressResponse {
  sessionId: string;
  progressPercent: number;
  lastActiveAt: string;
}

interface ReviewSummaryResponse {
  sessionId: string;
  overallThoughts?: string;
  recommendability?: string;
  notifyNewVersion: boolean;
  submittedAt: string;
}

export const readerApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getReaderLink: build.query<ReaderLinkResponse, string>({
      query: (token) => `/r/${token}`,
    }),
    startReview: build.mutation<
      StartReviewResponse,
      { token: string; readerName?: string; anonymous?: boolean }
    >({
      query: ({ token, ...body }) => ({
        url: `/r/${token}/start`,
        method: "POST",
        body,
      }),
    }),
    updateReviewProgress: build.mutation<
      ReviewProgressResponse,
      {
        sessionId: string;
        sectionId: string;
        paragraphId?: string;
        completed?: boolean;
        progressPercent?: number;
      }
    >({
      query: ({ sessionId, ...body }) => ({
        url: `/reviews/${sessionId}/progress`,
        method: "POST",
        body,
      }),
    }),
    addReviewReaction: build.mutation<
      Reaction,
      {
        sessionId: string;
        sectionId: string;
        paragraphId: string;
        type: ReactionType;
        text: string;
      }
    >({
      query: ({ sessionId, ...body }) => ({
        url: `/reviews/${sessionId}/reactions`,
        method: "POST",
        body,
      }),
    }),
    submitReviewSummary: build.mutation<
      ReviewSummaryResponse,
      {
        sessionId: string;
        overallThoughts?: string;
        recommendability?: "maybe" | "yes" | "absolutely";
        notifyNewVersion?: boolean;
      }
    >({
      query: ({ sessionId, ...body }) => ({
        url: `/reviews/${sessionId}/summary`,
        method: "POST",
        body,
      }),
    }),
  }),
});

export const {
  useGetReaderLinkQuery,
  useStartReviewMutation,
  useUpdateReviewProgressMutation,
  useAddReviewReactionMutation,
  useSubmitReviewSummaryMutation,
} = readerApi;
