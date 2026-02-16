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
const LoginRoute = lazy(() => import("@/routes/login"));
const RouteError = lazy(() => import("@/routes/route-error"));
const ProfileRoute = lazy(() => import("@/routes/profile"));
const PortfolioRoute = lazy(() => import("@/routes/portfolio"));
const OrcidCallbackRoute = lazy(() => import("@/routes/orcid-callback"));
const ZenodoCallbackRoute = lazy(() => import("@/routes/zenodo-callback"));
const SurfconextCallbackRoute = lazy(
  () => import("@/routes/surfconext-callback"),
);
const ProtectedRoute = lazy(() => import("@/routes/protected-route"));
const InboxRoute = lazy(() => import("@/routes/inbox"));
const ActivityRoute = lazy(() => import("@/routes/activity"));
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

const ProjectLayoutWrapper = lazy(() =>
  import("@/routes/project-layout-wrapper"),
);
const ProjectOverview = lazy(() =>
  import("@/features/project/pages/ProjectOverview").then((m) => ({
    default: m.ProjectOverview,
  })),
);
const ProjectTeamTab = lazy(() =>
  import("@/features/project/pages/ProjectTeamTab").then((m) => ({
    default: m.ProjectTeamTab,
  })),
);
const ProjectProductsTab = lazy(() =>
  import("@/features/project/pages/ProjectProductsTab").then((m) => ({
    default: m.ProjectProductsTab,
  })),
);
const ProjectActivityTab = lazy(() =>
  import("@/features/project/pages/ProjectActivityTab").then((m) => ({
    default: m.ProjectActivityTab,
  })),
);
const ProjectSettingsTab = lazy(() =>
  import("@/features/project/pages/ProjectSettingsTab").then((m) => ({
    default: m.ProjectSettingsTab,
  })),
);

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
              path: "activity",
              element: <ActivityRoute />,
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
              path: "portfolio",
              element: <PortfolioRoute />,
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
              path: ":id",
              element: <ProjectLayoutWrapper />,
              children: [
                {
                  index: true,
                  element: <ProjectOverview />,
                },
                {
                  path: "team",
                  element: <ProjectTeamTab />,
                },
                {
                  path: "products",
                  element: <ProjectProductsTab />,
                },
                {
                  path: "activity",
                  element: <ProjectActivityTab />,
                },
                {
                  path: "settings",
                  element: <ProjectSettingsTab />,
                },
                {
                  // Redirect old /edit route to /settings for backward compatibility
                  path: "edit",
                  element: <Navigate to="settings" replace />,
                },
              ],
            },
          ],
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
