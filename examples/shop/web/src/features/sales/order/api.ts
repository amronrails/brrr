import { api } from "@/lib/api/client";
import type { Order, OrderInput } from "./types";

export interface OrderList {
  data: Order[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/sales/orders";

export function listOrders(params?: { limit?: number; offset?: number }): Promise<OrderList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<OrderList>(qs ? `${base}?${qs}` : base);
}

export function getOrder(id: string): Promise<Order> {
  return api<Order>(`${base}/${id}`);
}

export function createOrder(input: OrderInput): Promise<Order> {
  return api<Order>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateOrder(id: string, input: OrderInput): Promise<Order> {
  return api<Order>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteOrder(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
