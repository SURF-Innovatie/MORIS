import { Folder } from "lucide-react";

import { ProjectEvent } from "@/api/events";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EventRenderer } from "@/components/events/EventRenderer";

interface ActivityHighlightsProps {
  eventHighlights: ProjectEvent[];
}

export const ActivityHighlights = ({
  eventHighlights,
}: ActivityHighlightsProps) => {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Folder className="h-5 w-5 text-primary" />
          Activity Highlights
        </CardTitle>
        <CardDescription>
          Recent updates across projects, roles, and publications.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {eventHighlights.length ? (
          <div className="space-y-4">
            {eventHighlights.map((event) => (
              <div
                key={event.id}
                className="rounded-lg border px-4 py-3"
              >
                <EventRenderer event={event} variant="compact" />
              </div>
            ))}
          </div>
        ) : (
          <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
            No recent activity to show yet.
          </div>
        )}
      </CardContent>
    </Card>
  );
};
