import { createContext, useContext, ReactNode, useMemo } from "react";
import { useParams } from "react-router-dom";
import { useGetProjectsIdAllowedEvents } from "@api/moris";
import { ProjectEventType } from "@/api/events";

interface UserEventsContextType {
  allowedEvents: string[];
  hasAccess: (event: ProjectEventType) => boolean;
  isLoading: boolean;
  isError: boolean;
}

const UserEventsContext = createContext<UserEventsContextType | undefined>(undefined);

export function UserEventsProvider({ children }: { children: ReactNode }) {
  const { id } = useParams<{ id: string }>();
  
  const { 
    data: allowedEvents, 
    isLoading, 
    isError 
  } = useGetProjectsIdAllowedEvents(id!, {
    query: {
      enabled: !!id,
    }
  });

  const value = useMemo(() => {
    const events = allowedEvents || [];
    return {
      allowedEvents: events,
      hasAccess: (event: ProjectEventType) => events.includes(event),
      isLoading,
      isError,
    };
  }, [allowedEvents, isLoading, isError]);

  return (
    <UserEventsContext.Provider value={value}>
      {children}
    </UserEventsContext.Provider>
  );
}

export function useUserEvents() {
  const context = useContext(UserEventsContext);
  if (context === undefined) {
    throw new Error("useUserEvents must be used within a UserEventsProvider");
  }
  return context;
}
