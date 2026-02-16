import { useParams } from "react-router-dom";
import { ChangelogTab } from "@/components/project-edit/ChangelogTab";

/**
 * ProjectActivityTab - Top-level tab for viewing project event timeline
 *
 * This component displays the project's activity history, showing all events
 * grouped by date. Events are rendered with context about what changed and who
 * made the change.
 */
export function ProjectActivityTab() {
  const { id } = useParams();

  if (!id) {
    return (
      <div className="flex h-64 items-center justify-center text-muted-foreground">
        Project ID is required
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <ChangelogTab projectId={id} />
    </div>
  );
}
