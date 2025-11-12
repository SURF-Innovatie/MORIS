import { isRouteErrorResponse, useRouteError } from 'react-router-dom';
import { AlertTriangle } from 'lucide-react';

import { Button } from '../components/ui/button';

const RouteError = () => {
  const error = useRouteError();

  const title = isRouteErrorResponse(error) ? `${error.status} - ${error.statusText}` : 'Something went wrong';
  const message = isRouteErrorResponse(error)
    ? error.data || 'We could not process your request.'
    : error instanceof Error
      ? error.message
      : 'An unexpected error occurred.';

  return (
    <div className="flex min-h-[60vh] flex-col items-center justify-center gap-4 text-center">
  <span className="flex h-16 w-16 items-center justify-center rounded-3xl bg-destructive/15 text-destructive">
  <AlertTriangle className="h-7 w-7" aria-hidden />
      </span>
      <div className="space-y-2">
        <h1 className="font-display text-3xl text-foreground">{title}</h1>
        <p className="text-sm text-muted-foreground">{message}</p>
      </div>
      <Button variant="ghost" onClick={() => window.location.replace('/')}>Return home</Button>
    </div>
  );
};

export default RouteError;
