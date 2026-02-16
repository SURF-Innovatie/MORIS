import { LucideIcon } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { StatusBadge } from "./status-badge";

export type ListItemVariant = "compact" | "default" | "detailed";

interface ListItemCommonProps {
  /**
   * Visual variant
   * @default "default"
   */
  variant?: ListItemVariant;

  /**
   * Primary text (title/name)
   */
  title: string;

  /**
   * Secondary text (email/description)
   */
  subtitle?: string;

  /**
   * Optional avatar image URL
   */
  avatarUrl?: string;

  /**
   * Optional avatar fallback text
   */
  avatarFallback?: string;

  /**
   * Optional icon (used instead of avatar)
   */
  icon?: LucideIcon;

  /**
   * Icon background color class (e.g., "bg-primary/10")
   */
  iconBgColor?: string;

  /**
   * Icon foreground color class (e.g., "text-primary")
   */
  iconFgColor?: string;

  /**
   * Badges to display
   */
  badges?: Array<{
    label: string;
    variant?: "default" | "secondary" | "outline" | "destructive";
    className?: string;
  }>;

  /**
   * Whether this item represents pending state
   */
  pending?: boolean;

  /**
   * Click handler
   */
  onClick?: () => void;

  /**
   * Optional action component (e.g., ActionMenu, buttons)
   */
  action?: React.ReactNode;

  /**
   * Additional content to display (detailed variant)
   */
  children?: React.ReactNode;

  /**
   * Additional CSS classes
   */
  className?: string;
}

/**
 * ListItem - Unified list item component
 *
 * Provides consistent list item patterns across the application with three variants:
 * - compact: Minimal, single-line items (sidebars, navigation)
 * - default: Standard list items with avatar and description
 * - detailed: Rich card-style items with additional content
 *
 * @example
 * ```tsx
 * // Compact variant (sidebar)
 * <ListItem
 *   variant="compact"
 *   title="Project Alpha"
 *   icon={Book}
 *   onClick={() => navigate('/project/1')}
 * />
 *
 * // Default variant (member list)
 * <ListItem
 *   title="John Doe"
 *   subtitle="john@example.com"
 *   avatarUrl="..."
 *   badges={[{ label: "Lead", variant: "secondary" }]}
 *   pending={true}
 *   action={<ActionMenu sections={...} />}
 * />
 *
 * // Detailed variant (project card)
 * <ListItem
 *   variant="detailed"
 *   title="Research Project"
 *   subtitle="Active project with 5 members"
 *   badges={[{ label: "Active" }]}
 *   onClick={() => navigate('/project/1')}
 * >
 *   <div className="mt-4">Additional project details...</div>
 * </ListItem>
 * ```
 */
export function ListItem({
  variant = "default",
  title,
  subtitle,
  avatarUrl,
  avatarFallback,
  icon: Icon,
  iconBgColor = "bg-primary/10",
  iconFgColor = "text-primary",
  badges = [],
  pending = false,
  onClick,
  action,
  children,
  className,
}: ListItemCommonProps) {
  // Compact variant
  if (variant === "compact") {
    return (
      <button
        onClick={onClick}
        className={cn(
          "w-full flex items-center gap-2 px-2 py-1.5 text-sm font-normal text-muted-foreground hover:text-foreground hover:bg-muted/50 rounded-md transition-colors truncate",
          pending && "opacity-70",
          className,
        )}
      >
        {Icon ? (
          <div
            className={cn(
              "flex items-center justify-center min-w-4 w-4 h-4 rounded-full",
              iconBgColor,
              iconFgColor,
            )}
          >
            <Icon className="h-2.5 w-2.5" />
          </div>
        ) : avatarUrl || avatarFallback ? (
          <Avatar className="h-5 w-5">
            <AvatarImage src={avatarUrl} />
            <AvatarFallback className="text-[10px]">
              {avatarFallback}
            </AvatarFallback>
          </Avatar>
        ) : null}
        <span className="truncate">{title}</span>
        {pending && <StatusBadge status="pending" size="xs" showIcon={false} />}
      </button>
    );
  }

  // Default variant
  if (variant === "default") {
    const Component = onClick ? "button" : "div";
    return (
      <Component
        onClick={onClick}
        className={cn(
          "w-full flex items-center justify-between rounded-lg border p-4 transition-colors",
          onClick && "hover:bg-muted/50 cursor-pointer",
          pending && "opacity-70",
          className,
        )}
      >
        <div className="flex items-center gap-4 min-w-0 flex-1">
          {Icon ? (
            <div
              className={cn(
                "flex items-center justify-center h-10 w-10 rounded-full shrink-0",
                iconBgColor,
                iconFgColor,
              )}
            >
              <Icon className="h-5 w-5" />
            </div>
          ) : (
            <Avatar className="h-10 w-10 shrink-0">
              <AvatarImage src={avatarUrl} />
              <AvatarFallback className="font-semibold text-primary">
                {avatarFallback}
              </AvatarFallback>
            </Avatar>
          )}
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <p className="font-semibold leading-none truncate">{title}</p>
              {pending && <StatusBadge status="pending" size="xs" />}
              {badges.map((badge, idx) => (
                <Badge
                  key={idx}
                  variant={badge.variant || "secondary"}
                  className={cn(
                    "text-[10px] h-5 px-1.5 font-normal capitalize",
                    badge.className,
                  )}
                >
                  {badge.label}
                </Badge>
              ))}
            </div>
            {subtitle && (
              <p className="text-sm text-muted-foreground mt-1 truncate">
                {subtitle}
              </p>
            )}
          </div>
        </div>
        {action && <div className="ml-2 shrink-0">{action}</div>}
      </Component>
    );
  }

  // Detailed variant
  const Component = onClick ? "button" : "div";
  return (
    <Component
      onClick={onClick}
      className={cn(
        "w-full rounded-lg border p-6 transition-all",
        onClick && "hover:shadow-md hover:border-primary/20 cursor-pointer",
        pending && "opacity-70",
        className,
      )}
    >
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center gap-4 min-w-0 flex-1">
          {Icon ? (
            <div
              className={cn(
                "flex items-center justify-center h-12 w-12 rounded-full shrink-0",
                iconBgColor,
                iconFgColor,
              )}
            >
              <Icon className="h-6 w-6" />
            </div>
          ) : (
            <Avatar className="h-12 w-12 shrink-0">
              <AvatarImage src={avatarUrl} />
              <AvatarFallback className="font-semibold text-primary">
                {avatarFallback}
              </AvatarFallback>
            </Avatar>
          )}
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <p className="font-semibold text-lg leading-none truncate">
                {title}
              </p>
              {pending && <StatusBadge status="pending" size="sm" />}
              {badges.map((badge, idx) => (
                <Badge
                  key={idx}
                  variant={badge.variant || "secondary"}
                  className={cn("text-xs h-6 px-2", badge.className)}
                >
                  {badge.label}
                </Badge>
              ))}
            </div>
            {subtitle && (
              <p className="text-sm text-muted-foreground mt-2">{subtitle}</p>
            )}
          </div>
        </div>
        {action && <div className="ml-2 shrink-0">{action}</div>}
      </div>
      {children}
    </Component>
  );
}
