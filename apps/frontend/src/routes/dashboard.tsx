import { useNavigate } from "react-router-dom";
import { Inbox } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { useNotifications } from "@/contexts/NotificationContext";
import { NotificationList } from "@/components/notifications/NotificationList";

const DashboardRoute = () => {
  const navigate = useNavigate();
  const { unreadCount } = useNotifications();

  return (
    <div className="max-w-5xl mx-auto space-y-8 px-4 py-8">
      {/* Welcome Section / Prompt */}
      <div className="flex flex-col md:flex-row gap-6 md:items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Stay updated with your projects and teams.
          </p>
        </div>

        <div className="flex gap-3">
          <Button onClick={() => navigate("/dashboard/projects/new")}>
            New Project
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Feed */}
        <div className="lg:col-span-2 space-y-6">
          {/* Inbox Preview (if unread) */}
          {unreadCount > 0 && (
            <Card className="bg-muted/30 border-dashed">
              <CardHeader className="pb-3">
                <CardTitle className="text-sm font-medium flex items-center gap-2">
                  <Inbox className="h-4 w-4 text-primary" />
                  <span>You have {unreadCount} unread notifications</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <Button
                  variant="link"
                  className="p-0 h-auto text-primary"
                  onClick={() => navigate("/dashboard/inbox")}
                >
                  View Inbox &rarr;
                </Button>
              </CardContent>
            </Card>
          )}

          <div className="space-y-4">
            <h3 className="text-lg font-semibold flex items-center gap-2">
              Recent Activity
            </h3>
            <Card>
              <CardContent className="p-0">
                <NotificationList limit={10} />
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Right Side - Explore / Suggestions */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="text-base">Explore</CardTitle>
              <CardDescription>Discover interesting projects.</CardDescription>
            </CardHeader>
            <CardContent className="text-sm text-muted-foreground">
              Coming soon...
            </CardContent>
          </Card>

          <Card className="bg-linear-to-br from-primary/5 to-transparent border-none shadow-none">
            <CardHeader>
              <CardTitle className="text-base">Pro Tip</CardTitle>
            </CardHeader>
            <CardContent className="text-sm">
              <p>
                Press{" "}
                <kbd className="pointer-events-none gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium text-muted-foreground opacity-100">
                  /
                </kbd>{" "}
                to search anywhere.
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default DashboardRoute;
