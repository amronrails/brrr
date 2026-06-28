import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useProject, useCreateProject, useUpdateProject } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  name: z.string(),
  key: z.string(),
  description: z.string(),
  archived: z.boolean(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  name: "",
  key: "",
  description: "",
  archived: false,
};

export function ProjectFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useProject(id);
  const create = useCreateProject();
  const update = useUpdateProject();

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
        key: data.key,
        description: data.description,
        archived: data.archived,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/projects/projects");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Project</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="name">Name</Label>
              <Input id="name" type="text" {...register("name")} />
              {errors.name && <p className="text-xs text-red-600">{errors.name.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="key">Key</Label>
              <Input id="key" type="text" {...register("key")} />
              {errors.key && <p className="text-xs text-red-600">{errors.key.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="description">Description</Label>
              <Input id="description" type="text" {...register("description")} />
              {errors.description && <p className="text-xs text-red-600">{errors.description.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="archived">Archived</Label>
              <input id="archived" type="checkbox" className="h-4 w-4" {...register("archived")} />
              {errors.archived && <p className="text-xs text-red-600">{errors.archived.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/projects/projects")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
