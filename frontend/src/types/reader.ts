export type InviteIdentityMode = "email" | "named" | "anonymous" | "preview";

export interface Reader {
  id: string;
  name: string;
  email?: string;
  avatar: string;
  articleId: string;
  progress: number;
  token: string;
  identityMode?: InviteIdentityMode;
  isPreview?: boolean;
  invitedAt: string;
  lastActiveAt?: string;
}

export interface InviteReaderDto {
  articleId: string;
  identityMode: InviteIdentityMode;
  displayName?: string;
  email?: string;
  note: string;
  isPreview?: boolean;
}
