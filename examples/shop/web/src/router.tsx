import { createBrowserRouter } from "react-router-dom";
import { ProtectedRoute, LoginPage, RegisterPage } from "@/features/auth";
import { DashboardLayout } from "@/components/layout/DashboardLayout";
import { DashboardPage } from "@/features/dashboard";
import { UsersPage } from "@/features/users";
import { CategoryListPage, CategoryFormPage } from "@/features/catalog/category";
import { ProductListPage, ProductFormPage } from "@/features/catalog/product";
import { OrderListPage, OrderFormPage } from "@/features/sales/order";
import { OrderItemListPage, OrderItemFormPage } from "@/features/sales/order_item";
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
          { path: "catalog/categories", element: <CategoryListPage /> },
          { path: "catalog/categories/new", element: <CategoryFormPage /> },
          { path: "catalog/categories/:id/edit", element: <CategoryFormPage /> },
          { path: "catalog/products", element: <ProductListPage /> },
          { path: "catalog/products/new", element: <ProductFormPage /> },
          { path: "catalog/products/:id/edit", element: <ProductFormPage /> },
          { path: "sales/orders", element: <OrderListPage /> },
          { path: "sales/orders/new", element: <OrderFormPage /> },
          { path: "sales/orders/:id/edit", element: <OrderFormPage /> },
          { path: "sales/order-items", element: <OrderItemListPage /> },
          { path: "sales/order-items/new", element: <OrderItemFormPage /> },
          { path: "sales/order-items/:id/edit", element: <OrderItemFormPage /> },
          // brrr:routes-fe
        ],
      },
    ],
  },
]);
