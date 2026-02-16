import { useNavigate, useParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useGetProjectsId } from "@api/moris";
import { ProjectOverview as ProjectOverviewComponent } from "@/components/projects/details/ProjectOverview";
import { ProjectTeamList } from "@/components/projects/details/ProjectTeamList";
import { ProjectProductList } from "@/components/projects/details/ProjectProductList";

export const ProjectOverview = () => {
  const { id } = useParams();
  const navigate = useNavigate();

  const {
    data: project,
    isLoading,
    error,
    refetch,
  } = useGetProjectsId(id!, {
    query: {
      enabled: !!id,
    },
  });

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex flex-col items-center justify-center gap-4 py-12">
        <p className="text-destructive">Failed to load project details.</p>
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          Back to Dashboard
        </Button>
      </div>
    );
  }

  // Render components similar to the original details page but in the new layout structure
  return (
    <div className="space-y-8">
      <ProjectOverviewComponent project={project} />
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <ProjectTeamList project={project} onRefresh={refetch} />
        <ProjectProductList project={project} />
      </div>
    </div>
  );
};
