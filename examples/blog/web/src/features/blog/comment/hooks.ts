import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { CommentInput } from "./types";

const key = ["blog", "comments"] as const;

export function useComments(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listComments(params),
  });
}

export function useComment(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getComment(id as string),
    enabled: !!id,
  });
}

export function useCreateComment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CommentInput) => api.createComment(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdateComment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: CommentInput }) => api.updateComment(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeleteComment() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteComment(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
