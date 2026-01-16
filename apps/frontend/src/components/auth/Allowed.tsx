import { ReactNode } from "react";
import { useAccess } from "@/context/AccessContext";

interface AllowedProps {
  event: string;
  children: ReactNode;
  fallback?: ReactNode;
}

export function Allowed({ event, children, fallback = null }: AllowedProps) {
  const { hasAccess } = useAccess();
  
  if (hasAccess(event)) {
    return <>{children}</>;
  }
  
  return <>{fallback}</>;
}
