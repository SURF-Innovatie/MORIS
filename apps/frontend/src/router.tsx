import { lazy, Suspense } from "react";
import type { RouteObject } from "react-router-dom";
import { createBrowserRouter, Navigate, Outlet } from "react-router-dom";

// Loading component for Suspense fallback
const PageLoader = () => (
  <div className="flex h-screen w-full items-center justify-center bg-background">
    <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
  </div>
);

// Lazy components
const RootLayout = lazy(() => import("@/routes/root"));
const DashboardRoute = lazy(() => import("@/routes/dashboard"));
const CreateProjectRoute = lazy(() => import("@/routes/project-form"));
const ProjectEditRoute = lazy(() => import("@/routes/project-edit"));
const ProjectDetailsRoute = lazy(() => import("@/routes/project-details"));
const LoginRoute = lazy(() => import("@/routes/login"));
const RouteError = lazy(() => import("@/routes/route-error"));
const ProfileRoute = lazy(() => import("@/routes/profile"));
const OrcidCallbackRoute = lazy(() => import("@/routes/orcid-callback"));
const ZenodoCallbackRoute = lazy(() => import("@/routes/zenodo-callback"));
const SurfconextCallbackRoute = lazy(
  () => import("@/routes/surfconext-callback"),
);
const ProtectedRoute = lazy(() => import("@/routes/protected-route"));
const InboxRoute = lazy(() => import("@/routes/inbox"));
const ProjectsRoute = lazy(() => import("@/routes/projects"));
const ProductsRoute = lazy(() => import("@/routes/products"));
const AdminUsersRoute = lazy(() => import("@/routes/admin-users"));
const AdminUserEditRoute = lazy(() => import("@/routes/admin-user-edit"));
const AdminOrganisationsRoute = lazy(() =>
  import("@/routes/admin-organisations").then((m) => ({
    default: m.AdminOrganisationsRoute,
  })),
);
const UserOrganisationsRoute = lazy(() =>
  import("@/routes/user-organisations").then((m) => ({
    default: m.UserOrganisationsRoute,
  })),
);
const UserOrganisationRolesRoute = lazy(() =>
  import("@/routes/user-organisation-roles").then((m) => ({
    default: m.UserOrganisationRolesRoute,
  })),
);
const MultiRoleManagementRoute = lazy(() =>
  import("@/routes/admin-organisation-roles").then((m) => ({
    default: m.MultiRoleManagementRoute,
  })),
);
const OrgAnalyticsRoute = lazy(() => import("@/routes/org-analytics"));

export function createAppRouter() {
  const routes: RouteObject[] = [
    {
      path: "/",
      element: (
        <Suspense fallback={<PageLoader />}>
          <LoginRoute />
        </Suspense>
      ),
      errorElement: <RouteError />,
    },
    {
      path: "/dashboard",
      element: (
        <Suspense fallback={<PageLoader />}>
          <ProtectedRoute />
        </Suspense>
      ),
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
              path: "settings",
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
              path: "organisations/:orgId/analytics",
              element: <OrgAnalyticsRoute />,
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
      element: (
        <Suspense fallback={<PageLoader />}>
          <ProtectedRoute />
        </Suspense>
      ),
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
      element: (
        <Suspense fallback={<PageLoader />}>
          <ProtectedRoute />
        </Suspense>
      ),
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
      element: (
        <Suspense fallback={<PageLoader />}>
          <ProtectedRoute />
        </Suspense>
      ),
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
      element: (
        <Suspense fallback={<PageLoader />}>
          <SurfconextCallbackRoute />
        </Suspense>
      ),
      errorElement: <RouteError />,
    },
    {
      path: "*",
      element: <Navigate to="/" replace />,
    },
  ];

  return createBrowserRouter(routes);
}
