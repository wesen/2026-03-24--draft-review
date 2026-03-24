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
    getReaders: build.query<Reader[], string>({
      query: (articleId) => `/articles/${articleId}/readers`,
      providesTags: ["Reader"],
    }),
    getReactions: build.query<Reaction[], string>({
      query: (articleId) => `/articles/${articleId}/reactions`,
      providesTags: ["Reaction"],
    }),
    inviteReader: build.mutation<Reader, { articleId: string; email: string; note: string }>({
      query: ({ articleId, ...body }) => ({
        url: `/articles/${articleId}/invite`,
        method: "POST",
        body,
      }),
      invalidatesTags: ["Reader"],
    }),
    addReaction: build.mutation<Reaction, { articleId: string; sectionId: string; paragraphId: string; type: string; text: string }>({
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
  useGetReadersQuery,
  useGetReactionsQuery,
  useInviteReaderMutation,
  useAddReactionMutation,
} = articleApi;
