import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { useComment, useCreateComment, useUpdateComment } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  body: z.string(),
  approved: z.boolean(),
  post_id: z.string().uuid(),
  author_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  body: "",
  approved: false,
  post_id: "",
  author_id: "",
};

export function CommentFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = useComment(id);
  const create = useCreateComment();
  const update = useUpdateComment();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<FormValues>({ resolver: zodResolver(schema), defaultValues: defaults });

  useEffect(() => {
    if (data) {
      reset({
        body: data.body,
        approved: data.approved,
        post_id: data.post_id,
        author_id: data.author_id,
      });
    }
  }, [data, reset]);

  const onSubmit = handleSubmit(async (values) => {
    if (isEdit && id) {
      await update.mutateAsync({ id, input: values });
    } else {
      await create.mutateAsync(values);
    }
    navigate("/blog/comments");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Comment</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="body">Body</Label>
              <Input id="body" type="text" {...register("body")} />
              {errors.body && <p className="text-xs text-red-600">{errors.body.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="approved">Approved</Label>
              <input id="approved" type="checkbox" className="h-4 w-4" {...register("approved")} />
              {errors.approved && <p className="text-xs text-red-600">{errors.approved.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="post_id">Post ID</Label>
              <Input id="post_id" type="text" {...register("post_id")} />
              {errors.post_id && <p className="text-xs text-red-600">{errors.post_id.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="author_id">Author ID</Label>
              <Input id="author_id" type="text" {...register("author_id")} />
              {errors.author_id && <p className="text-xs text-red-600">{errors.author_id.message}</p>}
            </div>
            <div className="flex gap-2">
              <Button type="submit" disabled={isSubmitting}>
                {isEdit ? "Save" : "Create"}
              </Button>
              <Button type="button" variant="outline" onClick={() => navigate("/blog/comments")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
