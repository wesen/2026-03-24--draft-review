import { baseApi } from "./baseApi";
import type { Article, Reaction, Reader } from "../types";

export const articleApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getArticles: build.query<Article[], void>({
      query: () => "/articles",
      providesTags: ["Article"],
    }),
    getArticle: build.query<Article, string>({
      query: (id) => `/articles/${id}`,
      providesTags: (_result, _err, id) => [{ type: "Article", id }],
    }),
    createArticle: build.mutation<Article, { title: string; author: string; intro: string }>({
      query: (body) => ({ url: "/articles", method: "POST", body }),
      invalidatesTags: ["Article"],
    }),
    updateArticle: build.mutation<Article, { id: string } & Partial<Article>>({
      query: ({ id, ...body }) => ({
        url: `/articles/${id}`,
        method: "PATCH",
        body,
      }),
      invalidatesTags: (_result, _err, { id }) => [
        { type: "Article", id },
        "Article",
      ],
    }),
    createVersion: build.mutation<
      Article,
      { articleId: string; label?: string; intro?: string; authorNote?: string }
    >({
      query: ({ articleId, ...body }) => ({
        url: `/articles/${articleId}/versions`,
        method: "POST",
        body,
      }),
      invalidatesTags: (_result, _err, { articleId }) => [
        { type: "Article", id: articleId },
        "Article",
      ],
    }),
    generateShareToken: build.mutation<{ token: string; url: string }, string>({
      query: (articleId) => ({
        url: `/articles/${articleId}/share-token`,
        method: "POST",
      }),
    }),
    exportArticle: build.mutation<
      { articleId: string; message: string; status: string },
      string
    >({
      query: (articleId) => ({
        url: `/articles/${articleId}/export`,
        method: "POST",
      }),
    }),
    getReaders: build.query<Reader[], string>({
      query: (articleId) => `/articles/${articleId}/readers`,
      providesTags: ["Reader"],
    }),
    getReactions: build.query<Reaction[], string>({
      query: (articleId) => `/articles/${articleId}/reactions`,
      providesTags: ["Reaction"],
    }),
    inviteReader: build.mutation<
      Reader,
      { articleId: string; email: string; note: string }
    >({
      query: ({ articleId, ...body }) => ({
        url: `/articles/${articleId}/invite`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Reader"],
    }),
    addReaction: build.mutation<
      Reaction,
      {
        articleId: string;
        sectionId: string;
        paragraphId: string;
        type: string;
        text: string;
      }
    >({
      query: ({ articleId, ...body }) => ({
        url: `/articles/${articleId}/reactions`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Reaction"],
    }),
  }),
});

export const {
  useGetArticlesQuery,
  useGetArticleQuery,
  useCreateArticleMutation,
  useUpdateArticleMutation,
  useCreateVersionMutation,
  useGenerateShareTokenMutation,
  useExportArticleMutation,
  useGetReadersQuery,
  useGetReactionsQuery,
  useInviteReaderMutation,
  useAddReactionMutation,
} = articleApi;
