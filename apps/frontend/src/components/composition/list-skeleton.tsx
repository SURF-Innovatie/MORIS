import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

export type ListSkeletonVariant = "compact" | "default" | "detailed";

interface ListSkeletonProps {
  /**
   * Visual variant to match ListItem
   * @default "default"
   */
  variant?: ListSkeletonVariant;

  /**
   * Number of skeleton items to render
   * @default 3
   */
  count?: number;

  /**
   * Additional CSS classes for container
   */
  className?: string;
}

const CompactSkeletonItem = () => (
  <div className="flex items-center gap-2 py-1.5 px-2">
    <Skeleton className="h-4 w-4 rounded-full" />
    <Skeleton className="h-4 flex-1" />
  </div>
);

const DefaultSkeletonItem = () => (
  <div className="flex items-center justify-between rounded-lg border p-4">
    <div className="flex items-center gap-4 flex-1">
      <Skeleton className="h-10 w-10 rounded-full" />
      <div className="flex-1 space-y-2">
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-3 w-1/2" />
      </div>
    </div>
    <Skeleton className="h-8 w-8" />
  </div>
);

const DetailedSkeletonItem = () => (
  <div className="rounded-lg border p-6">
    <div className="flex items-start justify-between mb-4">
      <div className="flex items-center gap-4 flex-1">
        <Skeleton className="h-12 w-12 rounded-full" />
        <div className="flex-1 space-y-2">
          <Skeleton className="h-5 w-3/4" />
          <Skeleton className="h-4 w-1/2" />
        </div>
      </div>
      <Skeleton className="h-8 w-8" />
    </div>
    <div className="space-y-2">
      <Skeleton className="h-3 w-full" />
      <Skeleton className="h-3 w-5/6" />
    </div>
  </div>
);

/**
 * ListSkeleton - Consistent loading state component
 *
 * Provides skeleton loaders that match ListItem variants for consistent
 * loading states across the application.
 *
 * @example
 * ```tsx
 * // Compact list loading
 * <ListSkeleton variant="compact" count={5} />
 *
 * // Default list loading
 * <ListSkeleton variant="default" count={3} />
 *
 * // Detailed card grid loading
 * <div className="grid gap-4 md:grid-cols-2">
 *   <ListSkeleton variant="detailed" count={4} />
 * </div>
 * ```
 */
export function ListSkeleton({
  variant = "default",
  count = 3,
  className,
}: ListSkeletonProps) {
  const SkeletonItem =
    variant === "compact"
      ? CompactSkeletonItem
      : variant === "detailed"
        ? DetailedSkeletonItem
        : DefaultSkeletonItem;

  const containerClass =
    variant === "compact"
      ? "space-y-0.5"
      : variant === "detailed"
        ? "space-y-4"
        : "space-y-4";

  return (
    <div className={cn(containerClass, className)}>
      {Array.from({ length: count }, (_, i) => (
        <SkeletonItem key={i} />
      ))}
    </div>
  );
}
