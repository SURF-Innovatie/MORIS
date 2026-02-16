import { useEffect } from "react";
import axios from "axios";

/**
 * Hook to track project access in the user's portfolio
 *
 * Automatically calls the track-access endpoint when a project is viewed.
 * This updates the recent_project_ids list in the user's portfolio.
 */
export function useTrackProjectAccess(projectId: string | undefined) {
  useEffect(() => {
    if (!projectId) return;

    // Call the track-access endpoint
    // Using axios directly to avoid generating a React Query hook for a simple POST
    const trackAccess = async () => {
      try {
        await axios.post(
          `/api/portfolio/track-access/${projectId}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
          }
        );
      } catch (error) {
        // Silently fail - tracking is not critical to the user experience
        console.debug("Failed to track project access:", error);
      }
    };

    // Small delay to avoid tracking immediately on navigation
    const timer = setTimeout(trackAccess, 1000);
    return () => clearTimeout(timer);
  }, [projectId]);
}
