import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useTask, useCreateTask, useUpdateTask } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  title: z.string(),
  description: z.string(),
  status: z.string(),
  priority: z.coerce.number().int(),
  due_date: z.string(),
  done: z.boolean(),
  project_id: z.string().uuid(),
  assignee_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  title: "",
  description: "",
  status: "",
  priority: 0,
  due_date: "",
  done: false,
  project_id: "",
  assignee_id: "",
};

export function TaskFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useTask(id);
  const create = useCreateTask();
  const update = useUpdateTask();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: defaults });

  useEffect(() => {
    if (data) {
      reset({
        title: data.title,
        description: data.description,
        status: data.status,
        priority: data.priority,
        due_date: data.due_date,
        done: data.done,
        project_id: data.project_id,
        assignee_id: data.assignee_id,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/projects/tasks");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Task</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="title">Title</Label>
              <Input id="title" type="text" {...register("title")} />
              {errors.title && <p className="text-xs text-red-600">{errors.title.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="description">Description</Label>
              <Input id="description" type="text" {...register("description")} />
              {errors.description && <p className="text-xs text-red-600">{errors.description.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="status">Status</Label>
              <Input id="status" type="text" {...register("status")} />
              {errors.status && <p className="text-xs text-red-600">{errors.status.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="priority">Priority</Label>
              <Input id="priority" type="number" step="any" {...register("priority", { valueAsNumber: true })} />
              {errors.priority && <p className="text-xs text-red-600">{errors.priority.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="due_date">Due Date</Label>
              <Input id="due_date" type="text" {...register("due_date")} />
              {errors.due_date && <p className="text-xs text-red-600">{errors.due_date.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="done">Done</Label>
              <input id="done" type="checkbox" className="h-4 w-4" {...register("done")} />
              {errors.done && <p className="text-xs text-red-600">{errors.done.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="project_id">Project ID</Label>
              <Input id="project_id" type="text" {...register("project_id")} />
              {errors.project_id && <p className="text-xs text-red-600">{errors.project_id.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="assignee_id">Assignee ID</Label>
              <Input id="assignee_id" type="text" {...register("assignee_id")} />
              {errors.assignee_id && <p className="text-xs text-red-600">{errors.assignee_id.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/projects/tasks")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
