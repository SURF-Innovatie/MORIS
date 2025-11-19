import { AlertCircle, RefreshCw } from 'lucide-react';
import { useBackendStatus } from '../../contexts/BackendStatusContext';
import { Button } from '../ui/button';

export const BackendOfflineAlert = () => {
  const { isOnline, lastError, checkHealth, isChecking } = useBackendStatus();

  if (isOnline) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm">
      <div className="mx-4 w-full max-w-md rounded-2xl border border-red-500/20 bg-card p-6 shadow-2xl">
        <div className="flex items-start gap-4">
          <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-red-500/10">
            <AlertCircle className="h-6 w-6 text-red-500" />
          </div>
          <div className="flex-1">
            <h2 className="text-lg font-semibold text-foreground">Backend Unavailable</h2>
            <p className="mt-2 text-sm text-muted-foreground">
              Unable to connect to the backend server. Please check your connection or try again later.
            </p>
            {lastError && (
              <div className="mt-3 rounded-lg bg-red-500/5 p-3">
                <p className="text-xs font-mono text-red-400">{lastError}</p>
              </div>
            )}
            <Button
              onClick={checkHealth}
              disabled={isChecking}
              className="mt-4 w-full"
              variant="outline"
            >
              {isChecking ? (
                <>
                  <RefreshCw className="mr-2 h-4 w-4 animate-spin" />
                  Checking...
                </>
              ) : (
                <>
                  <RefreshCw className="mr-2 h-4 w-4" />
                  Retry Connection
                </>
              )}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
};
