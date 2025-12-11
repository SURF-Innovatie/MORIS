import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "@/hooks/useAuth";

interface ProtectedRouteProps {
  children: React.ReactNode;
  requireSysAdmin?: boolean;
}

export const ProtectedRoute = ({
  children,
  requireSysAdmin,
}: ProtectedRouteProps) => {
  const { isAuthenticated, isLoading, user } = useAuth();

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-lg">Loading...</div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requireSysAdmin) {
    // Check for is_sys_admin flag.
    // Note: The structure depends on generated types. Assuming user.user.is_sys_admin based on previous context.
    const isSysAdmin = user?.is_sys_admin;

    if (!isSysAdmin) {
      return (
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-lg text-red-600">Access Denied: Insufficient permissions</div>
        </div>
      );
    }
  }

  return <>{children}</>;
};
