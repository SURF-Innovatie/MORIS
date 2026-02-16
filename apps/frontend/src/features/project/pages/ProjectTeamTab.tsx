import { useMemo } from "react";
import { useParams } from "react-router-dom";
import { Loader2 } from "lucide-react";

import { useGetProjectsId, useGetProjectsIdPendingEvents } from "@api/moris";
import { PeopleTab } from "@/components/project-edit/PeopleTab";
import { applyPendingEvents } from "@/lib/events/projection";

/**
 * ProjectTeamTab - Top-level tab for viewing and managing project team members
 *
 * This component displays team members with their roles. Users with appropriate
 * permissions can add/remove members and edit roles through the PeopleTab component.
 */
export function ProjectTeamTab() {
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
      <PeopleTab
        projectId={id!}
        members={projectedProject.members || []}
        onRefresh={refetchProject}
      />
    </div>
  );
}
