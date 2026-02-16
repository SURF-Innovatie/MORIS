import { useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { Book, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useGetPortfolioMe, useGetProjects } from "@api/moris";
import { ListItem, ListSkeleton } from "@/components/composition";

export const CompactProjectList = () => {
  const navigate = useNavigate();
  const { data: portfolio, isLoading: isLoadingPortfolio } =
    useGetPortfolioMe();
  const { data: projects, isLoading: isLoadingProjects } = useGetProjects();

  const recentProjects = useMemo(() => {
    if (
      !portfolio?.recent_project_ids?.length ||
      !projects?.length
    )
      return [];
    const projectMap = new Map(
      projects
        .filter((project) => project.id)
        .map((project) => [project.id!, project]),
    );
    return portfolio.recent_project_ids
      .map((id) => projectMap.get(id))
      .filter((project): project is NonNullable<typeof project> => !!project)
      .slice(0, 7);
  }, [portfolio, projects]);

  const isLoading = isLoadingPortfolio || isLoadingProjects;

  if (isLoading) {
    return (
      <div className="px-2">
        <ListSkeleton variant="compact" count={5} />
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center justify-between px-2 group">
        <h4 className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
          Recent
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
        {recentProjects.length === 0 ? (
          <p className="px-2 text-xs text-muted-foreground italic">
            No recent projects.
          </p>
        ) : (
          recentProjects.map((project) => (
            <ListItem
              key={project.id}
              variant="compact"
              title={project.title || "Untitled Project"}
              icon={Book}
              onClick={() => navigate(`/projects/${project.id}`)}
            />
          ))
        )}

        {portfolio?.recent_project_ids && portfolio.recent_project_ids.length > 7 && (
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
