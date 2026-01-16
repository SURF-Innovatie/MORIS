import type { RouteObject } from "react-router-dom";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";

import RootLayout from "@/routes/root";
import DashboardRoute from "@/routes/dashboard";
import CreateProjectRoute from "@/routes/project-form";

import ProjectEditRoute from "@/routes/project-edit";
import ProjectDetailsRoute from "@/routes/project-details";
import LoginRoute from "@/routes/login";
import RouteError from "@/routes/route-error";
import ProfileRoute from "@/routes/profile";
import OrcidCallbackRoute from "@/routes/orcid-callback";
import ZenodoCallbackRoute from "@/routes/zenodo-callback";
import SurfconextCallbackRoute from "@/routes/surfconext-callback";
import ProtectedRoute from "@/routes/protected-route";
import InboxRoute from "@/routes/inbox";
import ProjectsRoute from "@/routes/projects";
import ProductsRoute from "@/routes/products";
import AdminUsersRoute from "@/routes/admin-users";
import AdminUserEditRoute from "@/routes/admin-user-edit";
import { AdminOrganisationsRoute } from "@/routes/admin-organisations";
import { UserOrganisationsRoute } from "@/routes/user-organisations";
import { UserOrganisationRolesRoute } from "@/routes/user-organisation-roles";
import { MultiRoleManagementRoute } from "@/routes/admin-organisation-roles";

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
              element: <CreateProjectRoute />,
            },
            {
              path: "products",
              element: <ProductsRoute />,
            },
            {
              path: "profile",
              element: <ProfileRoute />,
            },
            {
              path: "organisations",
              element: <UserOrganisationsRoute />,
            },
            {
              path: "organisations/:nodeId/members",
              element: <UserOrganisationRolesRoute />,
            },
            {
              path: "admin",
              element: (
                <ProtectedRoute requireSysAdmin>
                  <Outlet />
                </ProtectedRoute>
              ),
              children: [
                {
                  path: "users",
                  element: <AdminUsersRoute />,
                },
                {
                  path: "users/:id/edit",
                  element: <AdminUserEditRoute />,
                },
                {
                  path: "organisations",
                  element: <AdminOrganisationsRoute />,
                },
                {
                  path: "organisations/:nodeId/roles",
                  element: <MultiRoleManagementRoute />,
                },
              ],
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
      path: "/zenodo-callback",
      element: <ProtectedRoute />,
      errorElement: <RouteError />,
      children: [
        {
          index: true,
          element: <ZenodoCallbackRoute />,
        },
      ],
    },
    {
      path: "/surfconext-callback",
      element: <SurfconextCallbackRoute />,
      errorElement: <RouteError />,
    },
    {
      path: "*",
      element: <Navigate to="/" replace />,
    },
  ];

  return createBrowserRouter(routes);
}
