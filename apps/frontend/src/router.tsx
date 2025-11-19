import type { RouteObject } from 'react-router-dom';
import { createBrowserRouter, Navigate } from 'react-router-dom';

import RootLayout from './routes/root';
import DashboardRoute from './routes/dashboard';
import LoginRoute from './routes/login';
import RegisterRoute from './routes/register';
import RouteError from './routes/route-error';
import ProtectedRoute from './routes/protected-route';

export function createAppRouter() {
  const routes: RouteObject[] = [
    {
      path: '/',
      element: <LoginRoute />,
      errorElement: <RouteError />,
    },
    {
      path: '/register',
      element: <RegisterRoute />,
      errorElement: <RouteError />,
    },
    {
      path: '/dashboard',
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
          ],
        },
      ],
    },
    {
      path: '*',
      element: <Navigate to="/" replace />,
    },
  ];

  return createBrowserRouter(routes);
}
