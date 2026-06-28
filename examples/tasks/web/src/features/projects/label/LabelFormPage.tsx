import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useLabel, useCreateLabel, useUpdateLabel } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  name: z.string(),
  color: z.string(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  name: "",
  color: "",
};

export function LabelFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useLabel(id);
  const create = useCreateLabel();
  const update = useUpdateLabel();

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
        color: data.color,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/projects/labels");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Label</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="name">Name</Label>
              <Input id="name" type="text" {...register("name")} />
              {errors.name && <p className="text-xs text-red-600">{errors.name.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="color">Color</Label>
              <Input id="color" type="text" {...register("color")} />
              {errors.color && <p className="text-xs text-red-600">{errors.color.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/projects/labels")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
