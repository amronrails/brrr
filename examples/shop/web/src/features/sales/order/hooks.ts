import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { OrderInput } from "./types";

const key = ["sales", "orders"] as const;

export function useOrders(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listOrders(params),
  });
}

export function useOrder(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getOrder(id as string),
    enabled: !!id,
  });
}

export function useCreateOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: OrderInput) => api.createOrder(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdateOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: OrderInput }) => api.updateOrder(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeleteOrder() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteOrder(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
