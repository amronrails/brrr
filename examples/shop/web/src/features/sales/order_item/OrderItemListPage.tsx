import { Link } from "react-router-dom";
import { useOrderItems, useDeleteOrderItem } from "./hooks";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

export function OrderItemListPage() {
  const { data, isLoading, isError } = useOrderItems();
  const remove = useDeleteOrderItem();

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">OrderItems</h1>
        <Link to="new">
          <Button>New OrderItem</Button>
        </Link>
      </div>

      <Card className="overflow-hidden">
        {isLoading ? (
          <div className="p-6 text-sm text-neutral-500">Loading…</div>
        ) : isError ? (
          <div className="p-6 text-sm text-red-600">Failed to load.</div>
        ) : (
          <table className="w-full text-left text-sm">
            <thead className="border-b border-neutral-200 bg-neutral-50 text-neutral-500">
              <tr>
                <th className="px-4 py-3 font-medium">Quantity</th>
                <th className="px-4 py-3 font-medium">Unit Price</th>
                <th className="px-4 py-3 font-medium">Order ID</th>
                <th className="px-4 py-3 font-medium">Product ID</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {data?.data.map((item) => (
                <tr key={item.id} className="border-b border-neutral-100 last:border-0">
                  <td className="px-4 py-3">{String(item.quantity)}</td>
                  <td className="px-4 py-3">{String(item.unit_price)}</td>
                  <td className="px-4 py-3">{String(item.order_id)}</td>
                  <td className="px-4 py-3">{String(item.product_id)}</td>
                  <td className="px-4 py-3">
                    <div className="flex justify-end gap-3">
                      <Link to={`${item.id}/edit`} className="text-neutral-700 hover:underline">
                        Edit
                      </Link>
                      <button
                        onClick={() => remove.mutate(item.id)}
                        className="text-red-600 hover:underline"
                      >
                        Delete
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>
    </div>
  );
}
