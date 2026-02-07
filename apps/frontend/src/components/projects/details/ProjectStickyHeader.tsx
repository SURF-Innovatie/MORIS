import { useNavigate } from "react-router-dom";
import { ArrowLeft, Pencil, Eye, LayoutTemplate, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ExportButton } from "@/components/projects/ExportButton";
import { useGetProjectsProjectIdPages } from "@api/moris";

interface ProjectStickyHeaderProps {
  projectId: string;
  title: string;
}

export function ProjectStickyHeader({
  projectId,
  title,
}: ProjectStickyHeaderProps) {
  const navigate = useNavigate();
  const { data: pages } = useGetProjectsProjectIdPages(projectId);
  const existingPage = pages && pages.length > 0 ? pages[0] : null;

  return (
    <header className="sticky top-0 z-10 border-b bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container flex h-16 items-center justify-between py-4">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/dashboard")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div className="flex flex-col">
            <h1 className="text-lg font-semibold leading-none tracking-tight">
              {title}
            </h1>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <ExportButton projectId={projectId} />

          {existingPage ? (
            <>
              <Button
                variant="outline"
                onClick={() => navigate(`/pages/${existingPage.slug}`)}
              >
                <Eye className="mr-2 h-4 w-4" />
                View Page
              </Button>
              <Button
                variant="outline"
                onClick={() => navigate(`/pages/${existingPage.slug}/edit`)}
              >
                <LayoutTemplate className="mr-2 h-4 w-4" />
                Edit Page
              </Button>
            </>
          ) : (
            <Button
              variant="outline"
              onClick={() =>
                navigate(`/pages/new?type=project&projectId=${projectId}`)
              }
            >
              <Plus className="mr-2 h-4 w-4" />
              Create Page
            </Button>
          )}

          <Button onClick={() => navigate(`/projects/${projectId}/edit`)}>
            <Pencil className="mr-2 h-4 w-4" />
            Edit Project
          </Button>
        </div>
      </div>
    </header>
  );
}
