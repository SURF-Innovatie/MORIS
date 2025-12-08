import { format } from "date-fns";
import { Link } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import { useGetUsersIdEventsApproved } from "@api/moris";

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
        <Card className="h-full border-dashed">
            <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
                <CardDescription>
                    Your recent publications and project updates.
                </CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col">
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
                                <div className="flex items-center justify-between">
                                    <Badge variant="outline" className="font-normal text-xs">
                                        {event.type}
                                    </Badge>
                                    <span className="text-xs text-muted-foreground">
                                        {event.at ? format(new Date(event.at), "PPP") : "N/A"}
                                    </span>
                                </div>
                                <p className="text-sm mt-1">{event.details}</p>
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
                            Once you start working on projects or publishing research,
                            your activity will appear here.
                        </p>
                    </div>
                )}
            </CardContent>
        </Card>
    );
}
