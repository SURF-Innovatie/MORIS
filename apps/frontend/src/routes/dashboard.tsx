import { useNavigate } from "react-router-dom";
import {
  Inbox,
  Bell,
  CheckCircle2,
} from "lucide-react";
import { format } from "date-fns";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
} from "@/components/ui/card";
import { useNotifications } from "@/context/NotificationContext";
import { ProjectList } from "@/components/projects/ProjectList";

const DashboardRoute = () => {
  const navigate = useNavigate();
  const { notifications, unreadCount, markAsRead } = useNotifications();

  const handleMarkAsRead = async (id: string) => {
    await markAsRead(id);
  };

  return (
    <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
      {/* Main Content (Projects) */}
      <div className="lg:col-span-2 space-y-8">
        <ProjectList />
      </div>

      {/* Sidebar (Inbox) */}
      <div className="space-y-8">
        <section>
          <div className="mb-6 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Inbox className="h-4 w-4 text-muted-foreground" />
              <h2 className="text-lg font-semibold tracking-tight">Inbox</h2>
              <Badge variant="secondary" className="ml-auto h-5 px-1.5 text-[10px]">
                {unreadCount}
              </Badge>
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="h-7 px-2 text-xs"
              onClick={() => navigate("/dashboard/inbox")}
            >
              View All
            </Button>
          </div>

          <Card className="border-none shadow-none bg-transparent">
            <CardContent className="p-0">
              {!notifications || notifications.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <Bell className="mb-2 h-8 w-8 text-muted-foreground/30" />
                  <p className="text-xs text-muted-foreground">
                    All caught up!
                  </p>
                </div>
              ) : (
                <div className="space-y-2">
                  {notifications.slice(0, 5).map((notification) => (
                    <div
                      key={notification.id}
                      className="group relative flex flex-col gap-1 rounded-lg border bg-card p-3 text-sm shadow-sm transition-all hover:shadow-md"
                    >
                      <div className="flex items-start justify-between gap-2">
                        <div className="font-medium leading-tight">
                          {!notification.read && (
                            <span className="mr-2 inline-block h-1.5 w-1.5 rounded-full bg-primary align-middle" />
                          )}
                          {notification.message}
                        </div>
                        {!notification.read && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 -mr-1 -mt-1 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity"
                            title="Mark as read"
                            onClick={() => handleMarkAsRead(notification.id!)}
                          >
                            <CheckCircle2 className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                      <p className="text-[10px] text-muted-foreground/70 mt-1">
                        {notification.sentAt ? format(new Date(notification.sentAt), "MMM d, h:mm a") : "Just now"}
                      </p>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </section>
      </div>
    </div>
  );
};

export default DashboardRoute;
