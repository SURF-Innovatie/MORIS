import { Home, Inbox, Activity } from "lucide-react";
import { useLocation, Link } from "react-router-dom";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { CompactProjectList } from "@/components/projects/CompactProjectList";
import { CompactOrganisationList } from "@/components/organisations/CompactOrganisationList";

export const AppNavigation = ({
  onItemClick,
}: {
  onItemClick?: () => void;
}) => {
  const location = useLocation();

  const mainNavItems = [
    { label: "Home", href: "/dashboard", icon: Home },
    { label: "Inbox", href: "/dashboard/inbox", icon: Inbox },
    { label: "Activity", href: "/dashboard/activity", icon: Activity },
  ];

  return (
    <div className="flex flex-col h-full py-4 space-y-6">
      {/* Main Links */}
      <div className="px-2 space-y-1">
        {mainNavItems.map((item) => {
          const Icon = item.icon;
          const isActive =
            location.pathname === item.href ||
            (item.href !== "/dashboard" &&
              location.pathname.startsWith(item.href));

          return (
            <Button
              key={item.href}
              variant="ghost"
              asChild
              className={cn(
                "w-full justify-start gap-3 px-3 py-2 h-9 font-normal text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors",
                isActive && "bg-muted text-foreground font-medium",
              )}
              onClick={onItemClick}
            >
              <Link to={item.href}>
                <Icon className="h-4 w-4" />
                {item.label}
              </Link>
            </Button>
          );
        })}
      </div>

      {/* Dividers & Lists */}
      <div className="px-2">
        <CompactProjectList />
      </div>

      <div className="px-2">
        <CompactOrganisationList />
      </div>
    </div>
  );
};
