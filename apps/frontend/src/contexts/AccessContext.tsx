import { createContext, useContext, ReactNode, useMemo } from "react";

interface AccessContextType {
  allowedEvents: string[];
  hasAccess: (event: string) => boolean;
  isLoading: boolean;
  isError: boolean;
}

const AccessContext = createContext<AccessContextType | undefined>(undefined);

interface AccessProviderProps {
  children: ReactNode;
  allowedEvents?: string[];
  isLoading?: boolean;
  isError?: boolean;
}

export function AccessProvider({
  children,
  allowedEvents = [],
  isLoading = false,
  isError = false
}: AccessProviderProps) {

  const value = useMemo(() => {
    return {
      allowedEvents,
      hasAccess: (event: string) => allowedEvents.includes(event),
      isLoading,
      isError,
    };
  }, [allowedEvents, isLoading, isError]);

  return (
    <AccessContext.Provider value={value}>
      {children}
    </AccessContext.Provider>
  );
}

export function useAccess() {
  const context = useContext(AccessContext);
  if (context === undefined) {
    throw new Error("useAccess must be used within an AccessProvider");
  }
  return context;
}
