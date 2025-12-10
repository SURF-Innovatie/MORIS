import type { RouteObject } from "react-router-dom";
import { createBrowserRouter, Navigate } from "react-router-dom";

import RootLayout from "./routes/root";
import DashboardRoute from "./routes/dashboard";
import ProjectFormRoute from "./routes/project-form";
import ProjectEditRoute from "./routes/project-edit";
import ProjectDetailsRoute from "./routes/project-details";
import LoginRoute from "./routes/login";
import RouteError from "./routes/route-error";
import ProfileRoute from "./routes/profile";
import OrcidCallbackRoute from "./routes/orcid-callback";
import ProtectedRoute from "./routes/protected-route";
import InboxRoute from "./routes/inbox";
import ProjectsRoute from "./routes/projects";
import ProductsRoute from "./routes/products";

export function createAppRouter() {
  const routes: RouteObject[] = [
    {
      path: "/",
      element: <LoginRoute />,
      errorElement: <RouteError />,
    },
    {
      path: "/dashboard",
      element: <ProtectedRoute />,
      errorElement: <RouteError />,
      children: [
        {
          element: <RootLayout />,
          children: [
            {
              index: true,
              element: <DashboardRoute />,
            },
            {
              path: "inbox",
              element: <InboxRoute />,
            },
            {
              path: "projects",
              element: <ProjectsRoute />,
            },
            {
              path: "projects/new",
              element: <ProjectFormRoute />,
            },
            {
              path: "products",
              element: <ProductsRoute />,
            },
            {
              path: "profile",
              element: <ProfileRoute />,
            },
          ],
        },
      ],
    },
    {
      path: "/projects",
      element: <ProtectedRoute />,
      errorElement: <RouteError />,
      children: [
        {
          path: ":id",
          element: <ProjectDetailsRoute />,
        },
        {
          path: ":id/edit",
          element: <ProjectEditRoute />,
        },
      ],
    },
    {
      path: "/orcid-callback",
      element: <ProtectedRoute />,
      errorElement: <RouteError />,
      children: [
        {
          index: true,
          element: <OrcidCallbackRoute />,
        },
      ],
    },
    {
      path: "*",
      element: <Navigate to="/" replace />,
    },
  ];

  return createBrowserRouter(routes);
}
