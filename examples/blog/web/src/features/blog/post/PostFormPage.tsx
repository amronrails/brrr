import { useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useNavigate, useParams } from "react-router-dom";
import { usePost, useCreatePost, useUpdatePost } from "./hooks";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const schema = z.object({
  title: z.string(),
  slug: z.string(),
  excerpt: z.string(),
  body: z.string(),
  published: z.boolean(),
  views: z.coerce.number().int(),
  author_id: z.string().uuid(),
});

type FormValues = z.infer<typeof schema>;

const defaults: FormValues = {
  title: "",
  slug: "",
  excerpt: "",
  body: "",
  published: false,
  views: 0,
  author_id: "",
};

export function PostFormPage() {
  const { id } = useParams();
  const isEdit = Boolean(id);
  const navigate = useNavigate();
  const { data } = usePost(id);
  const create = useCreatePost();
  const update = useUpdatePost();

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
        slug: data.slug,
        excerpt: data.excerpt,
        body: data.body,
        published: data.published,
        views: data.views,
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
    navigate("/blog/posts");
  });

  return (
    <div className="max-w-lg">
      <Card>
        <CardHeader>
          <CardTitle>{isEdit ? "Edit" : "New"} Post</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4" noValidate>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="title">Title</Label>
              <Input id="title" type="text" {...register("title")} />
              {errors.title && <p className="text-xs text-red-600">{errors.title.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="slug">Slug</Label>
              <Input id="slug" type="text" {...register("slug")} />
              {errors.slug && <p className="text-xs text-red-600">{errors.slug.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="excerpt">Excerpt</Label>
              <Input id="excerpt" type="text" {...register("excerpt")} />
              {errors.excerpt && <p className="text-xs text-red-600">{errors.excerpt.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="body">Body</Label>
              <Input id="body" type="text" {...register("body")} />
              {errors.body && <p className="text-xs text-red-600">{errors.body.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="published">Published</Label>
              <input id="published" type="checkbox" className="h-4 w-4" {...register("published")} />
              {errors.published && <p className="text-xs text-red-600">{errors.published.message}</p>}
            </div>
            <div className="flex flex-col gap-1.5">
              <Label htmlFor="views">Views</Label>
              <Input id="views" type="number" step="any" {...register("views", { valueAsNumber: true })} />
              {errors.views && <p className="text-xs text-red-600">{errors.views.message}</p>}
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
              <Button type="button" variant="outline" onClick={() => navigate("/blog/posts")}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
