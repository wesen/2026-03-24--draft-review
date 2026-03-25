import { baseApi } from "./baseApi";

export interface AuthUserInfo {
  authenticated: boolean;
  authMode: "dev" | "oidc";
  issuer?: string;
  subject?: string;
  email?: string;
  emailVerified?: boolean;
  preferredUsername?: string;
  displayName?: string;
  picture?: string;
  scopes?: string[];
  sessionExpiresAt?: string;
}

interface MeResponse {
  data: AuthUserInfo;
}

export const authApi = baseApi.injectEndpoints({
  endpoints: (build) => ({
    getMe: build.query<MeResponse, void>({
      query: () => "/me",
    }),
  }),
});

export const { useGetMeQuery } = authApi;
