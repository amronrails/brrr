export interface Post {
  id: string;
  title: string;
  slug: string;
  excerpt: string;
  body: string;
  published: boolean;
  views: number;
  author_id: string;
  created_at: string;
  updated_at: string;
}

export interface PostInput {
  title: string;
  slug: string;
  excerpt: string;
  body: string;
  published: boolean;
  views: number;
  author_id: string;
}
