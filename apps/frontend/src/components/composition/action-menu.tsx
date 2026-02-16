import { MoreHorizontal, LucideIcon } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";

export interface ActionMenuItem {
  /**
   * Menu item label
   */
  label: string;

  /**
   * Optional icon to display before label
   */
  icon?: LucideIcon;

  /**
   * Click handler
   */
  onClick: () => void;

  /**
   * Whether this is a destructive action (red text)
   */
  destructive?: boolean;

  /**
   * Whether this item is disabled
   */
  disabled?: boolean;
}

export interface ActionMenuSection {
  /**
   * Optional section label
   */
  label?: string;

  /**
   * Menu items in this section
   */
  items: ActionMenuItem[];
}

interface ActionMenuProps {
  /**
   * Menu sections (array of items or array of sections)
   */
  sections: ActionMenuSection[] | ActionMenuItem[];

  /**
   * Alignment of the dropdown
   * @default "end"
   */
  align?: "start" | "center" | "end";

  /**
   * Size of the trigger button
   * @default "icon"
   */
  size?: "default" | "sm" | "icon";

  /**
   * Additional CSS classes for trigger button
   */
  className?: string;

  /**
   * Optional custom trigger element (replaces default MoreHorizontal button)
   */
  trigger?: React.ReactNode;
}

/**
 * ActionMenu - Unified dropdown menu component
 *
 * Provides consistent action menus across the application with standardized
 * MoreHorizontal icon, sizing, and structure.
 *
 * @example
 * ```tsx
 * // Simple menu
 * <ActionMenu
 *   sections={[
 *     { label: "Edit", onClick: handleEdit },
 *     { label: "Delete", onClick: handleDelete, destructive: true }
 *   ]}
 * />
 *
 * // Sectioned menu
 * <ActionMenu
 *   sections={[
 *     {
 *       label: "Actions",
 *       items: [
 *         { label: "Edit", icon: Edit, onClick: handleEdit },
 *         { label: "Share", icon: Share, onClick: handleShare }
 *       ]
 *     },
 *     {
 *       items: [
 *         { label: "Delete", icon: Trash, onClick: handleDelete, destructive: true }
 *       ]
 *     }
 *   ]}
 * />
 * ```
 */
export function ActionMenu({
  sections,
  align = "end",
  size = "icon",
  className,
  trigger,
}: ActionMenuProps) {
  // Normalize input to always be sections
  const normalizedSections: ActionMenuSection[] = Array.isArray(sections)
    ? sections[0] && "items" in sections[0]
      ? (sections as ActionMenuSection[])
      : [{ items: sections as ActionMenuItem[] }]
    : [];

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        {trigger || (
          <Button
            variant="ghost"
            size={size}
            className={cn("h-8 w-8", className)}
          >
            <MoreHorizontal className="h-4 w-4" />
          </Button>
        )}
      </DropdownMenuTrigger>
      <DropdownMenuContent align={align}>
        {normalizedSections.map((section, sectionIdx) => (
          <div key={sectionIdx}>
            {section.label && <DropdownMenuLabel>{section.label}</DropdownMenuLabel>}
            {section.items.map((item, itemIdx) => {
              const Icon = item.icon;
              return (
                <DropdownMenuItem
                  key={itemIdx}
                  onClick={item.onClick}
                  disabled={item.disabled}
                  className={cn(item.destructive && "text-destructive")}
                >
                  {Icon && <Icon className="mr-2 h-4 w-4" />}
                  {item.label}
                </DropdownMenuItem>
              );
            })}
            {sectionIdx < normalizedSections.length - 1 && <DropdownMenuSeparator />}
          </div>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
