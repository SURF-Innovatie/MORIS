import { AlertCircle, CheckCircle, XCircle, Clock } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { getBadgeSizeClasses } from "@/lib/design-tokens";

export type BadgeStatus = "pending" | "approved" | "rejected" | "info";
export type BadgeSize = "xs" | "sm" | "md";

interface StatusBadgeProps {
  /**
   * The status to display
   */
  status: BadgeStatus;

  /**
   * Size variant
   * @default "xs"
   */
  size?: BadgeSize;

  /**
   * Whether to show an icon
   * @default true
   */
  showIcon?: boolean;

  /**
   * Optional custom label (overrides default status label)
   */
  label?: string;

  /**
   * Additional CSS classes
   */
  className?: string;
}

const statusConfig = {
  pending: {
    label: "Pending",
    icon: Clock,
    className: "border-yellow-500 text-yellow-600 bg-yellow-50",
  },
  approved: {
    label: "Approved",
    icon: CheckCircle,
    className: "border-green-500 text-green-600 bg-green-50",
  },
  rejected: {
    label: "Rejected",
    icon: XCircle,
    className: "border-red-500 text-red-600 bg-red-50",
  },
  info: {
    label: "Info",
    icon: AlertCircle,
    className: "border-blue-500 text-blue-600 bg-blue-50",
  },
} as const;

const iconSizeMap = {
  xs: "h-3 w-3",
  sm: "h-3.5 w-3.5",
  md: "h-4 w-4",
} as const;

/**
 * StatusBadge - Standardized status indicator component
 *
 * Provides consistent visual representation of status across the application.
 * Commonly used for pending events, approval states, and other status indicators.
 *
 * @example
 * ```tsx
 * <StatusBadge status="pending" size="xs" />
 * <StatusBadge status="approved" showIcon={false} />
 * <StatusBadge status="rejected" label="Denied" />
 * ```
 */
export function StatusBadge({
  status,
  size = "xs",
  showIcon = true,
  label,
  className,
}: StatusBadgeProps) {
  const config = statusConfig[status];
  const Icon = config.icon;
  const iconSize = iconSizeMap[size];
  const sizeClasses = getBadgeSizeClasses(size);

  return (
    <Badge
      variant="outline"
      className={cn(
        sizeClasses,
        config.className,
        "font-normal capitalize inline-flex items-center gap-1",
        className,
      )}
    >
      {showIcon && <Icon className={iconSize} />}
      <span>{label || config.label}</span>
    </Badge>
  );
}
