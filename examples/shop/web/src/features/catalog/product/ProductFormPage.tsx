import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useProduct, useCreateProduct, useUpdateProduct } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  name: z.string(),
  sku: z.string(),
  price: z.string(),
  stock: z.coerce.number().int(),
  active: z.boolean(),
  metadata: z.unknown(),
  category_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  name: "",
  sku: "",
  price: "",
  stock: 0,
  active: false,
  metadata: "",
  category_id: "",
};

export function ProductFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useProduct(id);
  const create = useCreateProduct();
  const update = useUpdateProduct();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: defaults });

  useEffect(() => {
    if (data) {
      reset({
        name: data.name,
        sku: data.sku,
        price: data.price,
        stock: data.stock,
        active: data.active,
        metadata: data.metadata,
        category_id: data.category_id,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/catalog/products");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Product</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="name">Name</Label>
              <Input id="name" type="text" {...register("name")} />
              {errors.name && <p className="text-xs text-red-600">{errors.name.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="sku">Sku</Label>
              <Input id="sku" type="text" {...register("sku")} />
              {errors.sku && <p className="text-xs text-red-600">{errors.sku.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="price">Price</Label>
              <Input id="price" type="text" {...register("price")} />
              {errors.price && <p className="text-xs text-red-600">{errors.price.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="stock">Stock</Label>
              <Input id="stock" type="number" step="any" {...register("stock", { valueAsNumber: true })} />
              {errors.stock && <p className="text-xs text-red-600">{errors.stock.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="active">Active</Label>
              <input id="active" type="checkbox" className="h-4 w-4" {...register("active")} />
              {errors.active && <p className="text-xs text-red-600">{errors.active.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="metadata">Metadata</Label>
              <Input id="metadata" type="text" {...register("metadata")} />
              {errors.metadata && <p className="text-xs text-red-600">{errors.metadata.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="category_id">Category ID</Label>
              <Input id="category_id" type="text" {...register("category_id")} />
              {errors.category_id && <p className="text-xs text-red-600">{errors.category_id.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/catalog/products")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
