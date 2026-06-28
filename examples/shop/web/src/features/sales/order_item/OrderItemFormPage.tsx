import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useOrderItem, useCreateOrderItem, useUpdateOrderItem } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  quantity: z.coerce.number().int(),
  unit_price: z.string(),
  order_id: z.string().uuid(),
  product_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  quantity: 0,
  unit_price: "",
  order_id: "",
  product_id: "",
};

export function OrderItemFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useOrderItem(id);
  const create = useCreateOrderItem();
  const update = useUpdateOrderItem();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: defaults });

  useEffect(() => {
    if (data) {
      reset({
        quantity: data.quantity,
        unit_price: data.unit_price,
        order_id: data.order_id,
        product_id: data.product_id,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/sales/order-items");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} OrderItem</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="quantity">Quantity</Label>
              <Input id="quantity" type="number" step="any" {...register("quantity", { valueAsNumber: true })} />
              {errors.quantity && <p className="text-xs text-red-600">{errors.quantity.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="unit_price">Unit Price</Label>
              <Input id="unit_price" type="text" {...register("unit_price")} />
              {errors.unit_price && <p className="text-xs text-red-600">{errors.unit_price.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="order_id">Order ID</Label>
              <Input id="order_id" type="text" {...register("order_id")} />
              {errors.order_id && <p className="text-xs text-red-600">{errors.order_id.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="product_id">Product ID</Label>
              <Input id="product_id" type="text" {...register("product_id")} />
              {errors.product_id && <p className="text-xs text-red-600">{errors.product_id.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/sales/order-items")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
