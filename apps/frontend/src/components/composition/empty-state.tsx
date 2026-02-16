import { LucideIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface EmptyStateProps {
  /**
   * Icon to display (Lucide icon component)
   */
  icon?: LucideIcon;

  /**
   * Title text
   */
  title: string;

  /**
   * Optional description text
   */
  description?: string;

  /**
   * Optional action button configuration
   */
  action?: {
    label: string;
    onClick: () => void;
  };

  /**
   * Size variant
   * @default "default"
   */
  size?: "sm" | "default" | "lg";

  /**
   * Additional CSS classes
   */
  className?: string;
}

const sizeConfig = {
  sm: {
    container: "py-6",
    icon: "h-8 w-8 mb-2",
    title: "text-sm",
    description: "text-xs",
  },
  default: {
    container: "py-8",
    icon: "h-10 w-10 mb-3",
    title: "text-base",
    description: "text-sm",
  },
  lg: {
    container: "py-12",
    icon: "h-12 w-12 mb-4",
    title: "text-lg",
    description: "text-base",
  },
} as const;

/**
 * EmptyState - Consistent empty state component
 *
 * Used across lists, tables, and content areas to indicate no data is available.
 * Supports optional icon, description, and action button.
 *
 * @example
 * ```tsx
 * <EmptyState
 *   icon={Building2}
 *   title="No projects found"
 *   description="Create your first project to get started"
 *   action={{ label: "Create Project", onClick: handleCreate }}
 * />
 * ```
 */
export function EmptyState({
  icon: Icon,
  title,
  description,
  action,
  size = "default",
  className,
}: EmptyStateProps) {
  const config = sizeConfig[size];

  return (
    <div
      className={cn(
        "flex flex-col items-center justify-center text-center text-muted-foreground",
        config.container,
        className,
      )}
    >
      {Icon && <Icon className={cn(config.icon, "opacity-20")} />}
      <p className={cn("font-medium text-foreground", config.title)}>{title}</p>
      {description && (
        <p className={cn("mt-1 text-muted-foreground", config.description)}>
          {description}
        </p>
      )}
      {action && (
        <Button className="mt-4" size="sm" onClick={action.onClick}>
          {action.label}
        </Button>
      )}
    </div>
  );
}
