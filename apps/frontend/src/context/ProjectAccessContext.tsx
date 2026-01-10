import { ReactNode } from "react";
import { useParams } from "react-router-dom";
import { useGetProjectsIdAllowedEvents } from "@api/moris";
import { AccessProvider } from "./AccessContext";

export function ProjectAccessProvider({ children }: { children: ReactNode }) {
  const { id } = useParams<{ id: string }>();

  const {
    data: allowedEvents,
    isLoading,
    isError,
  } = useGetProjectsIdAllowedEvents(id!, {
    query: {
      enabled: !!id,
      // Refetch on every navigation to ensure fresh permissions
      refetchOnMount: "always",
      // Also refetch when window regains focus
      refetchOnWindowFocus: true,
      // Don't use stale data - always show loading state while refetching
      staleTime: 0,
    },
  });

  return (
    <AccessProvider
      allowedEvents={allowedEvents}
      isLoading={isLoading}
      isError={isError}
    >
      {children}
    </AccessProvider>
  );
}
