import { Outlet, NavLink, useParams, Link } from "react-router-dom";
import { Star, Eye, BookOpen, Settings } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

// Mock hook if it doesn't exist yet, to be replaced
const useProjectMock = (id: string) => {
  return {
    data: {
      id,
      title: "Project Name", // This would come from API in reality
      owner: { name: "Owner Name" },
      isPrivate: true,
    },
    isLoading: false,
  };
};

export const ProjectLayout = () => {
  const { id } = useParams<{ id: string }>();
  // In a real implementation we would fetch project data here to display title/owner
  // const { data: project, isLoading } = useProject(id!);
  const { data: project } = useProjectMock(id!); // Placeholder

  if (!id) return null;

  return (
    <div className="flex flex-col min-h-[calc(100vh-3.5rem)]">
      {/* Project Header */}
      <div className="border-b bg-muted/20 px-4 pt-4 pb-0 md:px-8 lg:px-12">
        <div className="mx-auto max-w-7xl space-y-4">
          {/* Breadcrumbs / Title Line */}
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="flex items-center gap-2 text-lg md:text-xl">
              <BookOpen className="h-5 w-5 text-muted-foreground" />
              <Link to={`/dashboard`} className="text-primary hover:underline">
                {project?.owner?.name || "Owner"}
              </Link>
              <span className="text-muted-foreground">/</span>
              <Link
                to={`/dashboard/projects/${id}`}
                className="font-semibold text-primary hover:underline"
              >
                {project?.title || "Project Name"}
              </Link>
              <Badge
                variant="outline"
                className="ml-2 rounded-full text-xs font-normal"
              >
                {project?.isPrivate ? "Private" : "Public"}
              </Badge>
            </div>

            {/* Actions */}
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm" className="h-8 gap-2">
                <Eye className="h-4 w-4" />
                Watch
                <span className="ml-1 inline-flex h-5 items-center justify-center rounded-full bg-muted px-1.5 text-xs">
                  3
                </span>
              </Button>
              <Button variant="outline" size="sm" className="h-8 gap-2">
                <Star className="h-4 w-4" />
                Star
                <span className="ml-1 inline-flex h-5 items-center justify-center rounded-full bg-muted px-1.5 text-xs">
                  12
                </span>
              </Button>
            </div>
          </div>

          {/* Tabs Navigation */}
          <nav className="flex items-center gap-1 overflow-x-auto -mb-px">
            <ProjectTab
              to={`/dashboard/projects/${id}`}
              end
              icon={BookOpen}
              label="Overview"
            />
            {/* 
              GitHub uses "Settings" for the project administration.
              Since our "Edit" page is essentially the administration page (General, People, Policies),
              we will label it "Settings" in the UI to match GitHub patterns, even if the internal route is /edit.
            */}
            <ProjectTab
              to={`/dashboard/projects/${id}/edit`}
              icon={Settings}
              label="Settings"
            />
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
