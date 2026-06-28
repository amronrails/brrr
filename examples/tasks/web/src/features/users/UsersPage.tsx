import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import type { User } from "@/features/auth/types";
import { Card } from "@/components/ui/card";

export function UsersPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["admin", "users"],
    queryFn: () => api<{ data: User[] }>("/admin/users"),
  });

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-semibold">Users</h1>

      <Card className="overflow-hidden">
        {isLoading ? (
          <div className="p-6 text-sm text-neutral-500">Loading…</div>
        ) : isError ? (
          <div className="p-6 text-sm text-red-600">Failed to load users.</div>
        ) : (
          <table className="w-full text-left text-sm">
            <thead className="border-b border-neutral-200 bg-neutral-50 text-neutral-500">
              <tr>
                <th className="px-4 py-3 font-medium">Name</th>
                <th className="px-4 py-3 font-medium">Email</th>
                <th className="px-4 py-3 font-medium">Role</th>
              </tr>
            </thead>
            <tbody>
              {data?.data.map((u) => (
                <tr key={u.id} className="border-b border-neutral-100 last:border-0">
                  <td className="px-4 py-3">{u.name || "—"}</td>
                  <td className="px-4 py-3">{u.email}</td>
                  <td className="px-4 py-3">{u.role}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </Card>
    </div>
  );
}
