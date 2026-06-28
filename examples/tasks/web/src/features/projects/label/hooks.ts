import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { LabelInput } from "./types";

const key = ["projects", "labels"] as const;

export function useLabels(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listLabels(params),
  });
}

export function useLabel(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getLabel(id as string),
    enabled: !!id,
  });
}

export function useCreateLabel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: LabelInput) => api.createLabel(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdateLabel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: LabelInput }) => api.updateLabel(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeleteLabel() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteLabel(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
