import { api } from "@/lib/api/client";
import type { OrderItem, OrderItemInput } from "./types";

export interface OrderItemList {
  data: OrderItem[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/sales/order-items";

export function listOrderItems(params?: { limit?: number; offset?: number }): Promise<OrderItemList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<OrderItemList>(qs ? `${base}?${qs}` : base);
}

export function getOrderItem(id: string): Promise<OrderItem> {
  return api<OrderItem>(`${base}/${id}`);
}

export function createOrderItem(input: OrderItemInput): Promise<OrderItem> {
  return api<OrderItem>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateOrderItem(id: string, input: OrderItemInput): Promise<OrderItem> {
  return api<OrderItem>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteOrderItem(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
