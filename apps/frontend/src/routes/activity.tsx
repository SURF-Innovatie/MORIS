import { useState, useMemo } from "react";
import { History } from "lucide-react";
import { format } from "date-fns";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { EmptyState } from "@/components/composition";
import { useNotifications } from "@/contexts/NotificationContext";

type ActivityFilter = "all" | "your-projects" | "your-orgs";

/**
 * Activity Page - Full activity feed with filters
 *
 * Displays a comprehensive feed of all events across projects and organizations
 * the user has access to, with filtering options.
 */
export default function ActivityRoute() {
  const [filter, setFilter] = useState<ActivityFilter>("all");
  const navigate = useNavigate();
  const { notifications } = useNotifications();

  // Filter out approval requests and apply user-selected filter
  const filteredActivities = useMemo(() => {
    if (!notifications) return [];

    let filtered = notifications.filter((n) => n.type !== "approval_request");

    // Apply filter (for now, just show all since we don't have org info)
    // TODO: Implement proper filtering when API provides org context
    return filtered;
  }, [notifications, filter]);

  const isLoading = false;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Activity</h1>
        <p className="text-muted-foreground mt-2">
          View all recent activity across your projects and organizations
        </p>
      </div>

      {/* Filters */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">Filter Activity</CardTitle>
          <CardDescription>
            Show activity from specific sources
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            <Button
              variant={filter === "all" ? "default" : "outline"}
              size="sm"
              onClick={() => setFilter("all")}
            >
              All Activity
            </Button>
            <Button
              variant={filter === "your-projects" ? "default" : "outline"}
              size="sm"
              onClick={() => setFilter("your-projects")}
            >
              Your Projects
            </Button>
            <Button
              variant={filter === "your-orgs" ? "default" : "outline"}
              size="sm"
              onClick={() => setFilter("your-orgs")}
            >
              Your Organizations
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Activity Feed */}
      <Card>
        <CardHeader>
          <CardTitle className="text-lg">
            {filter === "all"
              ? "All Activity"
              : filter === "your-projects"
                ? "Your Projects Activity"
                : "Your Organizations Activity"}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="flex justify-center py-8">
              <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
            </div>
          ) : filteredActivities.length === 0 ? (
            <EmptyState
              icon={History}
              title="No activity yet"
              description={
                filter === "all"
                  ? "Activity from your projects and organizations will appear here"
                  : filter === "your-projects"
                    ? "Activity from your projects will appear here"
                    : "Activity from your organizations will appear here"
              }
            />
          ) : (
            <div className="space-y-3">
              {filteredActivities.map((notification) => (
                <div
                  key={notification.id}
                  className="flex items-start gap-4 p-4 rounded-lg border hover:bg-muted/30 transition-colors cursor-pointer"
                  onClick={() => {
                    if (notification.project_id) {
                      navigate(`/projects/${notification.project_id}/activity`);
                    }
                  }}
                >
                  <div className="flex-1 space-y-1">
                    <p className="text-sm font-medium">{notification.message}</p>
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      {notification.sent_at && (
                        <span>
                          {format(
                            new Date(notification.sent_at),
                            "EEEE, MMMM d, yyyy 'at' h:mm a",
                          )}
                        </span>
                      )}
                      {!notification.read && (
                        <>
                          <span>•</span>
                          <span className="text-primary font-medium">
                            Unread
                          </span>
                        </>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
