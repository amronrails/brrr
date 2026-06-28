import { api } from "@/lib/api/client";
import type { Task, TaskInput } from "./types";

export interface TaskList {
  data: Task[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/projects/tasks";

export function listTasks(params?: { limit?: number; offset?: number }): Promise<TaskList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<TaskList>(qs ? `${base}?${qs}` : base);
}

export function getTask(id: string): Promise<Task> {
  return api<Task>(`${base}/${id}`);
}

export function createTask(input: TaskInput): Promise<Task> {
  return api<Task>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateTask(id: string, input: TaskInput): Promise<Task> {
  return api<Task>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteTask(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
