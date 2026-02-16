import { useState } from "react";
import { Bell, CheckCircle, ClipboardCheck } from "lucide-react";
import { format } from "date-fns";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { EmptyState, ListItem } from "@/components/composition";
import { useNotifications } from "@/contexts/NotificationContext";
import { ApprovalModal } from "@/components/project-edit/ApprovalModal";

/**
 * ActionableItems - Dashboard component showing pending approvals and notifications
 *
 * Displays items that require user action, such as:
 * - Pending event approvals
 * - Review requests
 * - Important notifications
 */
export function ActionableItems() {
  const { notifications } = useNotifications();
  const [selectedApproval, setSelectedApproval] = useState<{
    projectId: string;
    eventId: string;
    message: string;
  } | null>(null);

  // Filter for approval requests only
  const approvalItems =
    notifications?.filter((n) => n.type === "approval_request" && !n.read) ||
    [];

  const isLoading = false;
  const pendingCount = approvalItems.length;

  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle className="flex items-center gap-2">
            <Bell className="h-5 w-5" />
            Actionable Items
            {pendingCount > 0 && (
              <Badge size="xs" variant="destructive">
                {pendingCount}
              </Badge>
            )}
          </CardTitle>
          <CardDescription>Items requiring your attention</CardDescription>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center py-8">
            <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
          </div>
        ) : approvalItems.length === 0 ? (
          <EmptyState
            icon={CheckCircle}
            title="All caught up!"
            description="No pending items requiring your action"
            size="sm"
          />
        ) : (
          <div className="space-y-3">
            {approvalItems.map((notification) => (
              <ListItem
                key={notification.id}
                variant="default"
                title="Approval Required"
                subtitle={`${notification.message} • ${notification.sent_at ? format(new Date(notification.sent_at), "MMM d 'at' h:mm a") : "Just now"}`}
                icon={ClipboardCheck}
                iconBgColor="bg-blue-500/10"
                iconFgColor="text-blue-500"
                badges={[
                  {
                    label: "Pending Approval",
                    variant: "destructive" as const,
                  },
                ]}
                onClick={() => {
                  if (notification.project_id && notification.event_id) {
                    setSelectedApproval({
                      projectId: notification.project_id,
                      eventId: notification.event_id,
                      message: notification.message || "No details available",
                    });
                  }
                }}
                className="cursor-pointer"
              />
            ))}
          </div>
        )}
      </CardContent>

      {selectedApproval && (
        <ApprovalModal
          isOpen={!!selectedApproval}
          onClose={() => setSelectedApproval(null)}
          projectId={selectedApproval.projectId}
          eventId={selectedApproval.eventId}
          notificationMessage={selectedApproval.message}
        />
      )}
    </Card>
  );
}
