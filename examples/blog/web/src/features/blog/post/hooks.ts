import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { PostInput } from "./types";

const key = ["blog", "posts"] as const;

export function usePosts(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listPosts(params),
  });
}

export function usePost(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getPost(id as string),
    enabled: !!id,
  });
}

export function useCreatePost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: PostInput) => api.createPost(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdatePost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: PostInput }) => api.updatePost(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeletePost() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deletePost(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
