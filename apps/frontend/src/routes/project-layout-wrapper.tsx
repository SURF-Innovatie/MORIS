import { ProjectAccessProvider } from "@/contexts/ProjectAccessContext";
import { ProjectLayout } from "@/components/layout";

/**
 * ProjectLayoutWrapper - Wraps ProjectLayout with ProjectAccessProvider
 *
 * This wrapper ensures that all project pages have access to permission checking
 * via the useAccess() hook, which depends on ProjectAccessProvider being present
 * in the component tree.
 */
export default function ProjectLayoutWrapper() {
  return (
    <ProjectAccessProvider>
      <ProjectLayout />
    </ProjectAccessProvider>
  );
}
