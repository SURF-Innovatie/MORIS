import { useMemo } from "react";
import { Pin, FolderKanban } from "lucide-react";
import { useNavigate } from "react-router-dom";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EmptyState, ListItem, ListSkeleton } from "@/components/composition";
import { useGetPortfolioMe, useGetProjects } from "@api/moris";

/**
 * PinnedProjects - Dashboard component showing user's pinned projects
 *
 * Displays projects the user has pinned from their portfolio for quick access.
 */
export function PinnedProjects() {
  const navigate = useNavigate();

  const { data: portfolio, isLoading: isLoadingPortfolio } =
    useGetPortfolioMe();
  const { data: projects, isLoading: isLoadingProjects } = useGetProjects();

  const pinnedProjects = useMemo(() => {
    if (
      !portfolio?.pinned_project_ids?.length ||
      !projects?.length
    )
      return [];
    const projectMap = new Map(
      projects
        .filter((project) => project.id)
        .map((project) => [project.id!, project]),
    );
    return portfolio.pinned_project_ids
      .map((id) => projectMap.get(id))
      .filter((project): project is NonNullable<typeof project> => !!project);
  }, [portfolio, projects]);

  const isLoading = isLoadingPortfolio || isLoadingProjects;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2 text-lg">
          <Pin className="h-4 w-4" />
          Pinned Projects
        </CardTitle>
        <CardDescription>Quick access to your favorites</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <ListSkeleton variant="compact" count={3} />
        ) : pinnedProjects.length === 0 ? (
          <EmptyState
            icon={FolderKanban}
            title="No pinned projects"
            description="Pin projects from your portfolio for quick access"
            size="sm"
            action={{
              label: "Go to Portfolio",
              onClick: () => navigate("/dashboard/portfolio"),
            }}
          />
        ) : (
          <div className="space-y-2">
            {pinnedProjects.map((project) => (
              <ListItem
                key={project.id}
                variant="compact"
                title={project.title || "Untitled Project"}
                icon={FolderKanban}
                onClick={() => navigate(`/projects/${project.id}`)}
              />
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
