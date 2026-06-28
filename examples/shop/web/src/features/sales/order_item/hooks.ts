import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as api from "./api";
import type { OrderItemInput } from "./types";

const key = ["sales", "order_items"] as const;

export function useOrderItems(params?: { limit?: number; offset?: number }) {
  return useQuery({
    queryKey: [...key, params ?? {}],
    queryFn: () => api.listOrderItems(params),
  });
}

export function useOrderItem(id: string | undefined) {
  return useQuery({
    queryKey: [...key, id],
    queryFn: () => api.getOrderItem(id as string),
    enabled: !!id,
  });
}

export function useCreateOrderItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: OrderItemInput) => api.createOrderItem(input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useUpdateOrderItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (vars: { id: string; input: OrderItemInput }) => api.updateOrderItem(vars.id, vars.input),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}

export function useDeleteOrderItem() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => api.deleteOrderItem(id),
    onSuccess: () => qc.invalidateQueries({ queryKey: key }),
  });
}
