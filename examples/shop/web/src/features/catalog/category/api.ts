import { api } from "@/lib/api/client";
import type { Category, CategoryInput } from "./types";

export interface CategoryList {
  data: Category[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/catalog/categories";

export function listCategories(params?: { limit?: number; offset?: number }): Promise<CategoryList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<CategoryList>(qs ? `${base}?${qs}` : base);
}

export function getCategory(id: string): Promise<Category> {
  return api<Category>(`${base}/${id}`);
}

export function createCategory(input: CategoryInput): Promise<Category> {
  return api<Category>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateCategory(id: string, input: CategoryInput): Promise<Category> {
  return api<Category>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteCategory(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
