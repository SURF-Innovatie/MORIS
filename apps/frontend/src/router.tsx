import type { RouteObject } from 'react-router-dom';
import { QueryClient } from '@tanstack/react-query';
import { createBrowserRouter } from 'react-router-dom';

import RootLayout from './routes/root';
import DashboardRoute from './routes/dashboard';
import UsersRoute, { loader as usersLoader } from './routes/users';
import RouteError from './routes/route-error';

export function createAppRouter(queryClient: QueryClient) {
  const routes: RouteObject[] = [
    {
      path: '/',
      element: <RootLayout />,
      errorElement: <RouteError />,
      children: [
        {
          index: true,
          element: <DashboardRoute />,
        },
        {
          path: 'users',
          loader: usersLoader(queryClient),
          element: <UsersRoute />,
        },
      ],
    },
  ];

  return createBrowserRouter(routes);
}
