import { useState } from "react";
import { format } from "date-fns";
import { Link } from "react-router-dom";
import { Bell, CheckCircle2, ClipboardCheck, Info } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useNotifications } from "@/contexts/NotificationContext";
import { ApprovalModal } from "@/components/project-edit/ApprovalModal";
import { EMPTY_UUID } from "@/lib/utils";

interface NotificationListProps {
  limit?: number;
}

enum NotificationType {
  ApprovalRequest = "approval_request",
  StatusUpdate = "status_update",
  Info = "info",
}

const getIcon = (type?: string) => {
  switch (type) {
    case NotificationType.ApprovalRequest:
      return <ClipboardCheck className="h-5 w-5 text-blue-500" />;
    case NotificationType.StatusUpdate:
      return <Info className="h-5 w-5 text-green-500" />;
    default:
      return <Bell className="h-5 w-5 text-gray-500" />;
  }
};

export function NotificationList({ limit }: NotificationListProps) {
  const { notifications, markAsRead } = useNotifications();
  const [selectedApproval, setSelectedApproval] = useState<{
    projectId: string;
    eventId: string;
    message: string;
  } | null>(null);

  const handleMarkAsRead = async (id: string, type?: string) => {
    // Approval requests are only marked as read when the decision is made (handled by backend)
    if (type === NotificationType.ApprovalRequest) {
      return;
    }
    await markAsRead(id);
  };

  if (!notifications || notifications.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <Bell className="mb-4 h-12 w-12 text-muted-foreground/30" />
        <p className="text-lg font-medium text-muted-foreground">
          All caught up!
        </p>
        <p className="text-sm text-muted-foreground/70">
          You have no new notifications.
        </p>
      </div>
    );
  }

  const displayedNotifications = limit
    ? notifications.slice(0, limit)
    : notifications;

  return (
    <>
      <div className="divide-y divide-border">
        {displayedNotifications.map((notification, index) => {
          const isApproval =
            notification.type === NotificationType.ApprovalRequest;

          const Content = (
            <div
              className={`group flex items-start gap-4 p-4 transition-colors hover:bg-muted/50 ${
                !notification.read ? "bg-primary/5" : ""
              } ${index === 0 ? "rounded-t-xl" : ""} ${index === notifications.length - 1 ? "rounded-b-xl" : ""}`}
            >
              <div className="mt-1">{getIcon(notification.type)}</div>
              <div className="flex-1 space-y-1">
                <p
                  className={`text-sm ${!notification.read ? "font-medium text-foreground" : "text-muted-foreground"}`}
                >
                  {notification.message}
                </p>
                <p className="text-xs text-muted-foreground/70">
                  {notification.sentAt
                    ? format(
                        new Date(notification.sentAt),
                        "MMM d, yyyy 'at' h:mm a"
                      )
                    : "Just now"}
                </p>
              </div>
              {!notification.read &&
                notification.type !== NotificationType.ApprovalRequest && (
                  <Button
                    variant="ghost"
                    size="icon"
                    className="h-8 w-8 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Mark as read"
                    onClick={(e) => {
                      e.preventDefault();
                      e.stopPropagation();
                      handleMarkAsRead(notification.id!, notification.type);
                    }}
                  >
                    <CheckCircle2 className="h-4 w-4" />
                  </Button>
                )}
            </div>
          );

          if (isApproval && notification.projectId && notification.eventId) {
            return (
              <div
                key={notification.id}
                onClick={() => {
                  // No manual mark as read on click for approval
                  setSelectedApproval({
                    projectId: notification.projectId!,
                    eventId: notification.eventId!,
                    message: notification.message || "No details available",
                  });
                }}
                className="cursor-pointer"
              >
                {Content}
              </div>
            );
          }

          const linkTarget =
            notification.projectId && notification.projectId != EMPTY_UUID
              ? `/projects/${notification.projectId}`
              : null;

          return (
            <div key={notification.id}>
              {linkTarget ? (
                <Link
                  to={linkTarget}
                  className="block"
                  onClick={() => {
                    if (!notification.read) {
                      handleMarkAsRead(notification.id!, notification.type);
                    }
                  }}
                >
                  {Content}
                </Link>
              ) : (
                Content
              )}
            </div>
          );
        })}
      </div>

      {selectedApproval && (
        <ApprovalModal
          isOpen={!!selectedApproval}
          onClose={() => setSelectedApproval(null)}
          projectId={selectedApproval.projectId}
          eventId={selectedApproval.eventId}
          notificationMessage={selectedApproval.message}
        />
      )}
    </>
  );
}
