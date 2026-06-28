import { api } from "@/lib/api/client";
import type { Product, ProductInput } from "./types";

export interface ProductList {
  data: Product[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/catalog/products";

export function listProducts(params?: { limit?: number; offset?: number }): Promise<ProductList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<ProductList>(qs ? `${base}?${qs}` : base);
}

export function getProduct(id: string): Promise<Product> {
  return api<Product>(`${base}/${id}`);
}

export function createProduct(input: ProductInput): Promise<Product> {
  return api<Product>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateProduct(id: string, input: ProductInput): Promise<Product> {
  return api<Product>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteProduct(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
