import { useAuth } from "@/features/auth/AuthContext";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export function DashboardPage() {
  const { user } = useAuth();

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold">Dashboard</h1>
        <p className="text-sm text-neutral-500">Welcome back, {user?.name || user?.email}.</p>
      </div>

      <Card className="max-w-md">
        <CardHeader>
          <CardTitle>You are signed in</CardTitle>
          <CardDescription>Role: {user?.role}</CardDescription>
        </CardHeader>
        <CardContent className="text-sm text-neutral-600">
          Generate CRUD modules with{" "}
          <code className="rounded bg-neutral-100 px-1 py-0.5">brrr generate</code> and they will
          show up in the sidebar.
        </CardContent>
      </Card>
    </div>
  );
}
