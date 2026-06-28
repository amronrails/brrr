import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useOrder, useCreateOrder, useUpdateOrder } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  status: z.string(),
  total: z.string(),
  placed_at: z.string().datetime(),
  customer_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  status: "",
  total: "",
  placed_at: "",
  customer_id: "",
};

export function OrderFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useOrder(id);
  const create = useCreateOrder();
  const update = useUpdateOrder();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: defaults });

  useEffect(() => {
    if (data) {
      reset({
        status: data.status,
        total: data.total,
        placed_at: data.placed_at,
        customer_id: data.customer_id,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/sales/orders");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Order</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="status">Status</Label>
              <Input id="status" type="text" {...register("status")} />
              {errors.status && <p className="text-xs text-red-600">{errors.status.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="total">Total</Label>
              <Input id="total" type="text" {...register("total")} />
              {errors.total && <p className="text-xs text-red-600">{errors.total.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="placed_at">Placed At</Label>
              <Input id="placed_at" type="text" {...register("placed_at")} />
              {errors.placed_at && <p className="text-xs text-red-600">{errors.placed_at.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="customer_id">Customer ID</Label>
              <Input id="customer_id" type="text" {...register("customer_id")} />
              {errors.customer_id && <p className="text-xs text-red-600">{errors.customer_id.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/sales/orders")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
