import { api } from "@/lib/api/client";
import type { Post, PostInput } from "./types";

export interface PostList {
  data: Post[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/blog/posts";

export function listPosts(params?: { limit?: number; offset?: number }): Promise<PostList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<PostList>(qs ? `${base}?${qs}` : base);
}

export function getPost(id: string): Promise<Post> {
  return api<Post>(`${base}/${id}`);
}

export function createPost(input: PostInput): Promise<Post> {
  return api<Post>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updatePost(id: string, input: PostInput): Promise<Post> {
  return api<Post>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deletePost(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
