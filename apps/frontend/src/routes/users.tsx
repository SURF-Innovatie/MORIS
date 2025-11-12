import { QueryClient } from '@tanstack/react-query';
import { LoaderFunctionArgs } from 'react-router-dom';
import { Loader2 as Spinner, Users as UsersIcon } from 'lucide-react';

import { getGetAdminUsersListQueryOptions, useGetAdminUsersList } from '../api/generated-orval/moris';
import { Button } from '../components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/card';
import { Badge } from '../components/ui/badge';

export const loader = (queryClient: QueryClient) => async (_args: LoaderFunctionArgs) => {
  const options = getGetAdminUsersListQueryOptions();
  await queryClient.ensureQueryData(options);
  return null;
};

const UsersRoute = () => {
  const {
    data,
    error,
    refetch,
    isFetching,
    isPending,
  } = useGetAdminUsersList({
    query: {
      retry: 0,
      staleTime: 1000 * 60,
    },
  });

  // The backend returns a JSON string representation of the list. Parse if present.
  let users: Array<{ id: number; name: string; email?: string }> = [];
  if (data) {
    try {
      const parsed = JSON.parse(data);
      users = parsed?.users ?? [];
    } catch (err) {
      // If parsing fails, fall back to an empty array
      users = [];
    }
  }

  return (
    <section className="flex flex-col gap-8">
      <header className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <Badge variant="outline" className="mb-3">
            Connected to Go backend
          </Badge>
          <h2 className="flex items-center gap-3 font-display text-3xl tracking-tight">
            <span className="inline-flex h-12 w-12 items-center justify-center rounded-2xl bg-primary/15 text-primary">
              <UsersIcon className="h-6 w-6" aria-hidden />
            </span>
            User directory
          </h2>
          <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
            This screen uses Orval generated hooks (`useGetAdminUsersList`) to retrieve data via Axios and expose it
            through TanStack Query. Trigger a refetch to see optimistic updates in action.
          </p>
        </div>
          <Button variant="secondary" onClick={() => refetch()} disabled={isFetching}>
          <Spinner className={`mr-2 h-4 w-4 ${isFetching ? 'animate-spin' : ''}`} aria-hidden />
          Refresh
        </Button>
      </header>

      <Card>
        <CardHeader>
          <CardTitle>Team members</CardTitle>
          <CardDescription>
            Server state {isFetching ? 'is syncing…' : 'is fresh'} — {users.length} entr{users.length === 1 ? 'y' : 'ies'} loaded
            from the Go API.
          </CardDescription>
        </CardHeader>
        <CardContent>
          {isPending ? (
            <ul className="flex flex-col gap-4">
              {[...Array(4)].map((_, index) => (
                <li key={index} className="animate-pulse rounded-xl border border-white/5 bg-white/5 p-4" />
              ))}
            </ul>
          ) : error ? (
            <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-6 text-destructive">
              <h3 className="font-semibold">Could not reach the backend</h3>
              <p className="mt-2 text-sm opacity-80">
                {typeof error === 'object' && error && 'message' in error
                  ? String((error as Record<string, unknown>).message)
                  : 'Please make sure the Go API is running on port 8080.'}
              </p>
            </div>
          ) : users.length === 0 ? (
            <div className="rounded-xl border border-dashed border-white/15 bg-white/5 p-10 text-center">
              <p className="text-sm text-muted-foreground">
                No users yet. Use the backend to seed data and refresh the view.
              </p>
            </div>
          ) : (
            <ul className="flex flex-col divide-y divide-white/5">
              {users.map((user) => (
                <li key={user.id} className="flex flex-wrap items-center justify-between gap-3 py-4">
                  <div>
                    <p className="text-base font-semibold text-foreground">{user.name}</p>
                    <p className="text-sm text-muted-foreground">{user.email}</p>
                  </div>
                  <Badge variant="success">ID: {user.id}</Badge>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </section>
  );
};

export default UsersRoute;
