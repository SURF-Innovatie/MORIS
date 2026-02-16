import { Home, Inbox, Activity, Briefcase, FolderKanban } from "lucide-react";
import { useLocation, Link } from "react-router-dom";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { CompactProjectList } from "@/components/projects/CompactProjectList";
import { CompactOrganisationList } from "@/components/organisations/CompactOrganisationList";
import { Badge } from "@/components/ui/badge";
import { useNotifications } from "@/contexts/NotificationContext";

/**
 * NavSection - Section header for navigation groups
 */
const NavSection = ({ title }: { title: string }) => (
  <h4 className="px-2 mb-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
    {title}
  </h4>
);

/**
 * NavItem - Navigation link item
 */
interface NavItemProps {
  label: string;
  href: string;
  icon: any;
  badge?: number;
  isActive?: boolean;
  onClick?: () => void;
}

const NavItem = ({ label, href, icon: Icon, badge, isActive, onClick }: NavItemProps) => (
  <Button
    variant="ghost"
    asChild
    className={cn(
      "w-full justify-start gap-3 px-3 py-2 h-9 font-normal text-muted-foreground hover:text-foreground hover:bg-muted/50 transition-colors",
      isActive && "bg-muted text-foreground font-medium",
    )}
    onClick={onClick}
  >
    <Link to={href}>
      <Icon className="h-4 w-4" />
      <span className="flex-1">{label}</span>
      {badge !== undefined && badge > 0 && (
        <Badge size="xs" className="ml-auto">
          {badge}
        </Badge>
      )}
    </Link>
  </Button>
);

export const AppNavigation = ({
  onItemClick,
}: {
  onItemClick?: () => void;
}) => {
  const location = useLocation();
  const { unreadCount } = useNotifications();

  const isActive = (href: string) => {
    if (href === "/dashboard") {
      return location.pathname === href;
    }
    return location.pathname.startsWith(href);
  };

  return (
    <div className="flex flex-col h-full py-4 space-y-6">
      {/* Primary Navigation */}
      <div className="space-y-1">
        <NavSection title="Primary Navigation" />
        <div className="px-2 space-y-0.5">
          <NavItem
            label="Dashboard"
            href="/dashboard"
            icon={Home}
            isActive={isActive("/dashboard")}
            onClick={onItemClick}
          />
          <NavItem
            label="Inbox"
            href="/dashboard/inbox"
            icon={Inbox}
            badge={unreadCount}
            isActive={isActive("/dashboard/inbox")}
            onClick={onItemClick}
          />
        </div>
      </div>

      {/* Your Work */}
      <div className="space-y-1">
        <NavSection title="Your Work" />
        <div className="px-2 space-y-0.5">
          <NavItem
            label="Your Portfolio"
            href="/dashboard/portfolio"
            icon={Briefcase}
            isActive={isActive("/dashboard/portfolio")}
            onClick={onItemClick}
          />
          <NavItem
            label="Your Projects"
            href="/dashboard/projects"
            icon={FolderKanban}
            isActive={isActive("/dashboard/projects")}
            onClick={onItemClick}
          />
          <NavItem
            label="Your Activity"
            href="/dashboard/activity"
            icon={Activity}
            isActive={isActive("/dashboard/activity")}
            onClick={onItemClick}
          />
        </div>
      </div>

      {/* Recent Projects */}
      <div>
        <CompactProjectList />
      </div>

      {/* Organizations */}
      <div>
        <CompactOrganisationList />
      </div>
    </div>
  );
};
