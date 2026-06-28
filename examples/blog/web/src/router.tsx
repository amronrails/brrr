import { createBrowserRouter } from "react-router-dom";
import { ProtectedRoute, LoginPage, RegisterPage } from "@/features/auth";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { DashboardPage } from "@/features/dashboard";
import { UsersPage } from "@/features/users";
import { PostListPage, PostFormPage } from "@/features/blog/post";
import { CommentListPage, CommentFormPage } from "@/features/blog/comment";
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
          { path: "blog/posts", element: <PostListPage /> },
          { path: "blog/posts/new", element: <PostFormPage /> },
          { path: "blog/posts/:id/edit", element: <PostFormPage /> },
          { path: "blog/comments", element: <CommentListPage /> },
          { path: "blog/comments/new", element: <CommentFormPage /> },
          { path: "blog/comments/:id/edit", element: <CommentFormPage /> },
          // brrr:routes-fe
        ],
      },
    ],
  },
]);
