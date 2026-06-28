import { Link } from "react-router-dom";
import { useTasks, useDeleteTask } from "./hooks";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

export function TaskListPage() {
  const { data, isLoading, isError } = useTasks();
  const remove = useDeleteTask();

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">Tasks</h1>
        <Link to="new">
          <Button>New Task</Button>
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
                <th className="px-4 py-3 font-medium">Title</th>
                <th className="px-4 py-3 font-medium">Description</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Priority</th>
                <th className="px-4 py-3 font-medium">Due Date</th>
                <th className="px-4 py-3 font-medium">Done</th>
                <th className="px-4 py-3 font-medium">Project ID</th>
                <th className="px-4 py-3 font-medium">Assignee ID</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {data?.data.map((item) => (
                <tr key={item.id} className="border-b border-neutral-100 last:border-0">
                  <td className="px-4 py-3">{String(item.title)}</td>
                  <td className="px-4 py-3">{String(item.description)}</td>
                  <td className="px-4 py-3">{String(item.status)}</td>
                  <td className="px-4 py-3">{String(item.priority)}</td>
                  <td className="px-4 py-3">{String(item.due_date)}</td>
                  <td className="px-4 py-3">{String(item.done)}</td>
                  <td className="px-4 py-3">{String(item.project_id)}</td>
                  <td className="px-4 py-3">{String(item.assignee_id)}</td>
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
