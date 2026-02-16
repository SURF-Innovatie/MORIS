import React from "react";
import { cn } from "@/lib/utils";
import { Link, useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { LucideIcon } from "lucide-react";

interface SidebarItem {
  icon?: LucideIcon;
  label: string;
  href: string;
  items?: SidebarItem[];
}

interface SidebarGroup {
  label?: string;
  items: SidebarItem[];
}

interface SidebarLayoutProps {
  children: React.ReactNode;
  sidebarGroups: SidebarGroup[];
}

export const SidebarLayout = ({
  children,
  sidebarGroups,
  extraSidebarContent,
}: SidebarLayoutProps & { extraSidebarContent?: React.ReactNode }) => {
  const location = useLocation();

  return (
    <div className="flex flex-1 flex-col md:flex-row">
      <aside className="w-full md:w-64 lg:w-72 shrink-0 border-r md:min-h-[calc(100vh-3.5rem)]">
        <div className="sticky top-14 h-[calc(100vh-3.5rem)] overflow-y-auto py-6 px-3">
          <nav className="space-y-6">
            {sidebarGroups.map((group, i) => (
              <div key={i} className="flex flex-col gap-2">
                {group.label && (
                  <h4 className="px-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-1">
                    {group.label}
                  </h4>
                )}
                <div className="space-y-1">
                  {group.items.map((item) => (
                    <SidebarNavLink
                      key={item.href}
                      item={item}
                      currentPath={location.pathname}
                    />
                  ))}
                </div>
              </div>
            ))}
            {extraSidebarContent && (
              <div className="pt-4 mt-4 border-t border-border">
                {extraSidebarContent}
              </div>
            )}
          </nav>
        </div>
      </aside>
      <div className="flex-1 px-4 py-6 md:px-8 lg:px-12 max-w-7xl mx-auto w-full">
        {children}
      </div>
    </div>
  );
};

const SidebarNavLink = ({
  item,
  currentPath,
}: {
  item: SidebarItem;
  currentPath: string;
}) => {
  const isActive =
    currentPath === item.href ||
    (item.href !== "/dashboard" && currentPath.startsWith(item.href)); // Basic active check
  const Icon = item.icon;

  return (
    <Button
      variant="ghost"
      asChild
      className={cn(
        "w-full justify-start gap-3 px-3 py-2 h-9 font-normal text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors",
        isActive && "bg-muted text-foreground font-medium",
      )}
    >
      <Link to={item.href}>
        {Icon && <Icon className="h-4 w-4" />}
        {item.label}
      </Link>
    </Button>
  );
};
