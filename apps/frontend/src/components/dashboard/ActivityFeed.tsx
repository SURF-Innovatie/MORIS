import { useState, useMemo } from "react";
import { History } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { format } from "date-fns";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { EmptyState } from "@/components/composition";
import { useNotifications } from "@/contexts/NotificationContext";

type ActivityFilter = "all" | "your-projects" | "your-orgs";

/**
 * ActivityFeed - Dashboard component showing recent activity
 *
 * Displays a filtered feed of recent events from projects and organizations
 * the user has access to.
 */
export function ActivityFeed() {
  const navigate = useNavigate();
  const [filter, setFilter] = useState<ActivityFilter>("all");
  const { notifications } = useNotifications();

  // Filter out approval requests (those go to ActionableItems)
  // and apply user-selected filter
  const filteredActivities = useMemo(() => {
    if (!notifications) return [];

    let filtered = notifications.filter((n) => n.type !== "approval_request");

    // Apply filter (for now, just show all since we don't have org info)
    // TODO: Implement proper filtering when API provides org context
    return filtered;
  }, [notifications, filter]);

  const isLoading = false;

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="flex items-center gap-2">
              <History className="h-5 w-5" />
              Recent Activity
            </CardTitle>
            <CardDescription>Latest updates from your work</CardDescription>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate("/dashboard/activity")}
          >
            View All
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {/* Filters */}
        <div className="flex flex-wrap gap-2 mb-4">
          <Button
            variant={filter === "all" ? "secondary" : "ghost"}
            size="xs"
            onClick={() => setFilter("all")}
          >
            All
          </Button>
          <Button
            variant={filter === "your-projects" ? "secondary" : "ghost"}
            size="xs"
            onClick={() => setFilter("your-projects")}
          >
            Projects
          </Button>
          <Button
            variant={filter === "your-orgs" ? "secondary" : "ghost"}
            size="xs"
            onClick={() => setFilter("your-orgs")}
          >
            Organizations
          </Button>
        </div>

        {/* Activity List */}
        {isLoading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
          </div>
        ) : filteredActivities.length === 0 ? (
          <EmptyState
            icon={History}
            title="No recent activity"
            description="Activity from your projects and organizations will appear here"
            size="sm"
          />
        ) : (
          <div className="space-y-4">
            {filteredActivities.slice(0, 5).map((notification) => (
              <div
                key={notification.id}
                className="flex flex-col gap-1 pb-3 border-b last:border-b-0 last:pb-0 cursor-pointer hover:bg-muted/30 rounded px-2 py-1 -mx-2 -my-1 transition-colors"
                onClick={() => {
                  if (notification.project_id) {
                    navigate(`/projects/${notification.project_id}/activity`);
                  }
                }}
              >
                <p className="text-sm font-medium">{notification.message}</p>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  {notification.sent_at && (
                    <span>
                      {format(new Date(notification.sent_at), "MMM d 'at' h:mm a")}
                    </span>
                  )}
                  {!notification.read && (
                    <>
                      <span>•</span>
                      <span className="text-primary font-medium">New</span>
                    </>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
