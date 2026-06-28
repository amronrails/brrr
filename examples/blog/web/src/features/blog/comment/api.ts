import { api } from "@/lib/api/client";
import type { Comment, CommentInput } from "./types";

export interface CommentList {
  data: Comment[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/blog/comments";

export function listComments(params?: { limit?: number; offset?: number }): Promise<CommentList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<CommentList>(qs ? `${base}?${qs}` : base);
}

export function getComment(id: string): Promise<Comment> {
  return api<Comment>(`${base}/${id}`);
}

export function createComment(input: CommentInput): Promise<Comment> {
  return api<Comment>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateComment(id: string, input: CommentInput): Promise<Comment> {
  return api<Comment>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteComment(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
