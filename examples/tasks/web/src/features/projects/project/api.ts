import { api } from "@/lib/api/client";
import type { Project, ProjectInput } from "./types";

export interface ProjectList {
  data: Project[];
  total: number;
  limit: number;
  offset: number;
}

const base = "/projects/projects";

export function listProjects(params?: { limit?: number; offset?: number }): Promise<ProjectList> {
  const q = new URLSearchParams();
  if (params?.limit != null) q.set("limit", String(params.limit));
  if (params?.offset != null) q.set("offset", String(params.offset));
  const qs = q.toString();
  return api<ProjectList>(qs ? `${base}?${qs}` : base);
}

export function getProject(id: string): Promise<Project> {
  return api<Project>(`${base}/${id}`);
}

export function createProject(input: ProjectInput): Promise<Project> {
  return api<Project>(base, { method: "POST", body: JSON.stringify(input) });
}

export function updateProject(id: string, input: ProjectInput): Promise<Project> {
  return api<Project>(`${base}/${id}`, { method: "PUT", body: JSON.stringify(input) });
}

export function deleteProject(id: string): Promise<void> {
  return api<void>(`${base}/${id}`, { method: "DELETE" });
}
