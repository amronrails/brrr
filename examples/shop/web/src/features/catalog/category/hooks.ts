import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { CategoryInput } from "./types";

const key = ["catalog", "categories"] as const;

export function useCategories(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listCategories(params),
  });
}

export function useCategory(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getCategory(id as string),
    enabled: !!id,
  });
}

export function useCreateCategory() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CategoryInput) => api.createCategory(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdateCategory() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: CategoryInput }) => api.updateCategory(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeleteCategory() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteCategory(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
