import { Outlet, NavLink, useParams, Link } from "react-router-dom";
import {
  BookOpen,
  Settings,
  Users,
  Package,
  Activity,
  Loader2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useGetProjectsId } from "@api/moris";
import { useAccess } from "@/contexts/AccessContext";
import { ProjectEventType } from "@/api/events";
import { useTrackProjectAccess } from "@/hooks/useTrackProjectAccess";

export const ProjectLayout = () => {
  const { id } = useParams<{ id: string }>();
  const { hasAccess } = useAccess();
  const canEdit = hasAccess(ProjectEventType.TitleChanged);

  // Track project access for recent projects list
  useTrackProjectAccess(id);

  const {
    data: project,
    isLoading,
    error,
  } = useGetProjectsId(id!, {
    query: {
      enabled: !!id,
    },
  });

  if (!id) return null;

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-destructive mb-4">Failed to load project</p>
          <Button variant="outline" onClick={() => (window.location.href = "/dashboard")}>
            Back to Dashboard
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col min-h-[calc(100vh-3.5rem)]">
      {/* Project Header */}
      <div className="border-b bg-muted/20 px-4 pt-4 pb-0 md:px-8 lg:px-12">
        <div className="mx-auto max-w-7xl space-y-4">
          {/* Breadcrumbs / Title Line */}
          <div className="flex items-center gap-2 text-lg md:text-xl">
            <BookOpen className="h-5 w-5 text-muted-foreground" />
            <Link to={`/dashboard`} className="text-primary hover:underline">
              {project?.owning_org_node?.name || "Organization"}
            </Link>
            <span className="text-muted-foreground">/</span>
            <Link
              to={`/projects/${id}`}
              className="font-semibold text-primary hover:underline"
            >
              {project?.title || "Project"}
            </Link>
          </div>

          {/* Tabs Navigation */}
          <nav className="flex items-center gap-1 overflow-x-auto -mb-px">
            <ProjectTab
              to={`/projects/${id}`}
              end
              icon={BookOpen}
              label="Overview"
            />
            <ProjectTab
              to={`/projects/${id}/team`}
              icon={Users}
              label="Team"
            />
            <ProjectTab
              to={`/projects/${id}/products`}
              icon={Package}
              label="Products"
            />
            <ProjectTab
              to={`/projects/${id}/activity`}
              icon={Activity}
              label="Activity"
            />
            {canEdit && (
              <ProjectTab
                to={`/projects/${id}/settings`}
                icon={Settings}
                label="Settings"
              />
            )}
          </nav>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 px-4 py-6 md:px-8 lg:px-12 max-w-7xl mx-auto w-full">
        <Outlet />
      </div>
    </div>
  );
};

const ProjectTab = ({
  to,
  icon: Icon,
  label,
  end,
}: {
  to: string;
  icon: any;
  label: string;
  end?: boolean;
}) => {
  return (
    <NavLink
      to={to}
      end={end}
      className={({ isActive }) =>
        cn(
          "flex items-center gap-2 border-b-2 border-transparent px-4 py-2.5 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors",
          isActive && "border-primary text-foreground",
        )
      }
    >
      <Icon className="h-4 w-4" />
      {label}
    </NavLink>
  );
};
