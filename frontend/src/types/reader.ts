export interface Reader {
  id: string;
  name: string;
  email: string;
  avatar: string;
  articleId: string;
  progress: number;
  token: string;
  invitedAt: string;
  lastActiveAt?: string;
}

export interface InviteReaderDto {
  articleId: string;
  email: string;
  note: string;
}
