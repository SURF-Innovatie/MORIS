import { useIsFetching, useIsMutating } from '@tanstack/react-query';
import { Spinner } from '@mynaui/icons-react';

import { cn } from '../../lib/utils';

const QueryIndicator = () => {
  const fetching = useIsFetching();
  const mutating = useIsMutating();
  const isBusy = fetching + mutating > 0;

  if (!isBusy) {
    return null;
  }

  return (
    <div
      className={cn(
        'fixed bottom-6 right-6 flex items-center gap-3 rounded-xl border border-white/10 bg-background/90 px-4 py-3 text-sm text-muted-foreground shadow-mynaui-sm backdrop-blur',
      )}
    >
      <Spinner className="h-4 w-4 animate-spin text-primary" aria-hidden />
      <span className="font-medium text-foreground">Syncing data…</span>
    </div>
  );
};

export default QueryIndicator;
