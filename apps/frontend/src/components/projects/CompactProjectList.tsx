import { useNavigate } from "react-router-dom";
import { Book, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useGetProjects } from "@api/moris";
import { Skeleton } from "@/components/ui/skeleton";

export const CompactProjectList = () => {
  const navigate = useNavigate();
  const { data: projects, isLoading } = useGetProjects();

  if (isLoading) {
    return (
      <div className="space-y-2 px-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
        <Skeleton className="h-4 w-5/6" />
      </div>
    );
  }

  // Take first 7 projects (simulating "Top/Recent")
  const topProjects = projects?.slice(0, 7) || [];

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center justify-between px-2 group">
        <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Top Repositories
        </h4>
        <Button
          variant="ghost"
          size="icon"
          className="h-5 w-5 opacity-0 group-hover:opacity-100 transition-opacity"
          onClick={() => navigate("/dashboard/projects/new")}
        >
          <Plus className="h-3 w-3" />
        </Button>
      </div>

      <div className="space-y-0.5">
        {topProjects.length === 0 ? (
          <p className="px-2 text-xs text-muted-foreground italic">
            No projects yet.
          </p>
        ) : (
          topProjects.map((project) => (
            <Button
              key={project.id}
              variant="ghost"
              className="w-full justify-start gap-2 px-2 py-1.5 h-auto text-sm font-normal text-muted-foreground hover:text-foreground truncate"
              onClick={() => navigate(`/projects/${project.id}/edit`)}
            >
              <div className="flex items-center justify-center min-w-4 w-4 h-4 rounded-full bg-primary/10 text-primary">
                <Book className="h-2.5 w-2.5" />
              </div>
              <span className="truncate">{project.title}</span>
            </Button>
          ))
        )}

        {projects && projects.length > 7 && (
          <Button
            variant="ghost"
            className="w-full justify-start gap-2 px-2 py-1.5 h-auto text-xs text-muted-foreground hover:text-foreground"
            onClick={() => navigate("/dashboard/projects")}
          >
            Show more...
          </Button>
        )}
      </div>
    </div>
  );
};
