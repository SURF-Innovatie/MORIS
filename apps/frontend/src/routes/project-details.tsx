import { useNavigate, useParams } from "react-router-dom";
import { Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useGetProjectsId } from "@api/moris";
import { ProjectStickyHeader } from "@/components/projects/details/ProjectStickyHeader";
import { ProjectOverview } from "@/components/projects/details/ProjectOverview";
import { ProjectTeamList } from "@/components/projects/details/ProjectTeamList";
import { ProjectProductList } from "@/components/projects/details/ProjectProductList";

export default function ProjectDetailsRoute() {
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
      <div className="flex h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-4">
        <p className="text-destructive">Failed to load project details.</p>
        <Button variant="outline" onClick={() => navigate("/dashboard")}>
          Back to Dashboard
        </Button>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <ProjectStickyHeader projectId={id!} title={project.title || "Project"} />

      <main className="container flex-1 py-8 space-y-8">
        <ProjectOverview project={project} />
        <ProjectTeamList project={project} onRefresh={refetch} />
        <ProjectProductList project={project} />
      </main>
    </div>
  );
}
