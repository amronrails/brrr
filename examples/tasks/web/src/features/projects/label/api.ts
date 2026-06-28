import { api } from "@/lib/api/client";
import type { Label, LabelInput } from "./types";

export interface LabelList {
  data: Label[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/projects/labels";

export function listLabels(params?: { limit?: number; offset?: number }): Promise<LabelList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<LabelList>(qs ? `${base}?${qs}` : base);
}

export function getLabel(id: string): Promise<Label> {
  return api<Label>(`${base}/${id}`);
}

export function createLabel(input: LabelInput): Promise<Label> {
  return api<Label>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateLabel(id: string, input: LabelInput): Promise<Label> {
  return api<Label>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteLabel(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
