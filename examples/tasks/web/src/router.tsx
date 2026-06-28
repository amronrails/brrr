import { createBrowserRouter } from "react-router-dom";
import { ProtectedRoute, LoginPage, RegisterPage } from "@/features/auth";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { DashboardPage } from "@/features/dashboard";
import { UsersPage } from "@/features/users";
import { ProjectListPage, ProjectFormPage } from "@/features/projects/project";
import { TaskListPage, TaskFormPage } from "@/features/projects/task";
import { LabelListPage, LabelFormPage } from "@/features/projects/label";
// brrr:imports-fe

export const router = createBrowserRouter([
  { path: "/login", element: <LoginPage /> },
  { path: "/register", element: <RegisterPage /> },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <DashboardLayout />,
        children: [
          { index: true, element: <DashboardPage /> },
          {
            element: <ProtectedRoute requireAdmin />,
            children: [{ path: "users", element: <UsersPage /> }],
          },
          { path: "projects/projects", element: <ProjectListPage /> },
          { path: "projects/projects/new", element: <ProjectFormPage /> },
          { path: "projects/projects/:id/edit", element: <ProjectFormPage /> },
          { path: "projects/tasks", element: <TaskListPage /> },
          { path: "projects/tasks/new", element: <TaskFormPage /> },
          { path: "projects/tasks/:id/edit", element: <TaskFormPage /> },
          { path: "projects/labels", element: <LabelListPage /> },
          { path: "projects/labels/new", element: <LabelFormPage /> },
          { path: "projects/labels/:id/edit", element: <LabelFormPage /> },
          // brrr:routes-fe
        ],
      },
    ],
  },
]);
