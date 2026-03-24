import { baseApi } from "./baseApi";
import type { ReaderArticle } from "../types";

interface ReaderLinkResponse {
  reader: { id: string; name: string };
  article: ReaderArticle;
}

export const readerApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getReaderLink: build.query<ReaderLinkResponse, string>({
      query: (token) => `/r/${token}`,
    }),
  }),
});

export const { useGetReaderLinkQuery } = readerApi;
