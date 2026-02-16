import { format } from "date-fns";
import { History, Loader2 } from "lucide-react";
import { useGetProjectsIdChangelog } from "@api/moris";

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EventRenderer } from "@/components/events/EventRenderer";
import { ProjectEvent } from "@/api/events";
import { EmptyState } from "@/components/composition";

interface ChangelogTabProps {
  projectId: string;
}

export function ChangelogTab({ projectId }: ChangelogTabProps) {
  const {
    data: changelogData,
    isLoading,
    error,
  } = useGetProjectsIdChangelog(projectId);

  if (isLoading) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="flex justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center text-sm text-destructive">
            Failed to load changelog
          </div>
        </CardContent>
      </Card>
    );
  }

  const events = changelogData?.events || [];

  // Group events by date
  const groupedEvents = events.reduce(
    (acc, event) => {
      const date = format(new Date(event.at!), "yyyy-MM-dd");
      if (!acc[date]) acc[date] = [];
      acc[date].push(event);
      return acc;
    },
    {} as Record<string, typeof events>,
  );

  return (
    <Card>
      <CardHeader>
        <CardTitle>Project History</CardTitle>
        <CardDescription>
          View the recent activity and changes in this project.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-8">
          {Object.entries(groupedEvents).map(([date, dateEvents]) => (
            <div key={date}>
              <h4 className="mb-4 text-sm font-medium text-muted-foreground sticky top-0 bg-background py-2">
                {format(new Date(date), "EEEE, MMMM d, yyyy")}
              </h4>
              <div className="space-y-6">
                {dateEvents.map((event) => (
                  <div
                    key={event.id}
                    className="flex flex-col gap-1 border-b pb-4 last:border-0 last:pb-0"
                  >
                    <EventRenderer
                      event={event as ProjectEvent}
                      variant="compact"
                    />
                  </div>
                ))}
              </div>
            </div>
          ))}
          {events.length === 0 && (
            <EmptyState
              icon={History}
              title="No history available"
              description="Project events will appear here once actions are taken"
            />
          )}
        </div>
      </CardContent>
    </Card>
  );
}
