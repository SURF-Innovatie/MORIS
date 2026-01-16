import { Link } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useGetUsersIdEventsApproved } from "@api/moris";
import { EventRenderer } from "@/components/events/EventRenderer";
import { ProjectEvent } from "@/api/events";

interface ProfileActivityProps {
  userId: string;
}

export function ProfileActivity({ userId }: ProfileActivityProps) {
  const { data: eventsData, isLoading: isLoadingEvents } =
    useGetUsersIdEventsApproved(userId, {
      query: {
        enabled: !!userId,
      },
    });

  return (
    <Card className="flex flex-col h-[calc(100vh-10rem)] border-dashed">
      <CardHeader>
        <CardTitle>Recent Activity</CardTitle>
        <CardDescription>
          Your recent publications and project updates.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex-1 overflow-y-auto min-h-0">
        {isLoadingEvents ? (
          <div className="py-12 text-center text-muted-foreground">
            Loading activity...
          </div>
        ) : eventsData?.events?.length ? (
          <div className="space-y-6">
            {eventsData.events.map((event) => (
              <div
                key={event.id}
                className="flex flex-col gap-1 border-b pb-4 last:border-0 last:pb-0"
              >
                <EventRenderer event={event as ProjectEvent} />
              </div>
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-12 text-center text-muted-foreground">
            <div className="h-12 w-12 rounded-full bg-muted/50 flex items-center justify-center mb-4">
              <Link className="h-6 w-6 opacity-20" />
            </div>
            <p className="font-medium">No recent activity</p>
            <p className="text-sm mt-1 max-w-xs mx-auto">
              Once you start working on projects or publishing research, your
              activity will appear here.
            </p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
