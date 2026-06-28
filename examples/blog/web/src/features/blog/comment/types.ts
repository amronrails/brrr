export interface Comment {
  id: string;
  body: string;
  approved: boolean;
  post_id: string;
  author_id: string;
  created_at: string;
  updated_at: string;
}

export interface CommentInput {
  body: string;
  approved: boolean;
  post_id: string;
  author_id: string;
}
