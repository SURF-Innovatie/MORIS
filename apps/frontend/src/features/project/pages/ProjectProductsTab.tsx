import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { Loader2 } from "lucide-react";

import { useGetProjectsId, useGetProjectsIdPendingEvents } from "@api/moris";
import { ProductsTab } from "@/components/project-edit/ProductsTab";
import { applyPendingEvents } from "@/lib/events/projection";

/**
 * ProjectProductsTab - Top-level tab for viewing and managing project products (research outputs)
 *
 * This component displays products associated with the project. Users with appropriate
 * permissions can add products via DOI import or Zenodo upload, and remove products.
 */
export function ProjectProductsTab() {
  const { id } = useParams();

  const {
    data: project,
    isLoading,
    refetch: refetchProject,
  } = useGetProjectsId(id!, {
    query: {
      enabled: !!id,
    },
  });

  const { data: pendingEventsData } = useGetProjectsIdPendingEvents(id!, {
    query: {
      enabled: !!id,
    },
  });

  // Apply pending events to show projected state
  const projectedProject = useMemo(() => {
    if (!project) return undefined;
    if (!pendingEventsData?.events) return project;
    return applyPendingEvents(project, pendingEventsData.events);
  }, [project, pendingEventsData]);

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!projectedProject) {
    return (
      <div className="flex h-64 items-center justify-center text-muted-foreground">
        Project not found
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <ProductsTab
        projectId={id!}
        products={projectedProject.products || []}
        onRefresh={refetchProject}
      />
    </div>
  );
}
