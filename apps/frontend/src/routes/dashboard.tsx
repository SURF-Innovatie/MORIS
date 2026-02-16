import { useNavigate } from "react-router-dom";
import { Inbox, Plus, FolderKanban } from "lucide-react";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { useNotifications } from "@/contexts/NotificationContext";
import { ActionableItems } from "@/components/dashboard/ActionableItems";
import { ActivityFeed } from "@/components/dashboard/ActivityFeed";
import { PinnedProjects } from "@/components/dashboard/PinnedProjects";

const DashboardRoute = () => {
  const navigate = useNavigate();
  const { unreadCount } = useNotifications();
  const { user } = useAuth();

  return (
    <div className="space-y-6">
      {/* Welcome Banner */}
      <div className="flex flex-col md:flex-row gap-4 md:items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">
            Welcome back{user?.name ? `, ${user.name.split(" ")[0]}` : ""}
          </h1>
          <p className="text-muted-foreground mt-1">
            Here's what's happening with your projects and teams
          </p>
        </div>

        <div className="flex gap-2">
          <Button onClick={() => navigate("/dashboard/projects/new")}>
            <Plus className="h-4 w-4 mr-2" />
            New Project
          </Button>
        </div>
      </div>

      {/* Quick Inbox Alert */}
      {unreadCount > 0 && (
        <Card className="bg-primary/5 border-primary/20">
          <CardContent className="flex items-center justify-between py-4">
            <div className="flex items-center gap-3">
              <Inbox className="h-5 w-5 text-primary" />
              <div>
                <p className="font-medium">
                  You have {unreadCount} unread notification
                  {unreadCount !== 1 ? "s" : ""}
                </p>
                <p className="text-sm text-muted-foreground">
                  Check your inbox for updates
                </p>
              </div>
            </div>
            <Button
              variant="outline"
              size="sm"
              onClick={() => navigate("/dashboard/inbox")}
            >
              View Inbox
            </Button>
          </CardContent>
        </Card>
      )}

      {/* Main Dashboard Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left Column (2/3) - Main Content */}
        <div className="lg:col-span-2 space-y-6">
          <ActionableItems />
          <ActivityFeed />
        </div>

        {/* Right Column (1/3) - Sidebar */}
        <div className="space-y-6">
          <PinnedProjects />

          {/* Quick Links */}
          <Card>
            <CardHeader>
              <CardTitle className="text-lg">Quick Links</CardTitle>
              <CardDescription>Navigate to common areas</CardDescription>
            </CardHeader>
            <CardContent className="space-y-2">
              <Button
                variant="ghost"
                className="w-full justify-start"
                size="sm"
                onClick={() => navigate("/dashboard/projects")}
              >
                <FolderKanban className="h-4 w-4 mr-2" />
                All Projects
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start"
                size="sm"
                onClick={() => navigate("/dashboard/portfolio")}
              >
                <Inbox className="h-4 w-4 mr-2" />
                Portfolio
              </Button>
              <Button
                variant="ghost"
                className="w-full justify-start"
                size="sm"
                onClick={() => navigate("/dashboard/activity")}
              >
                <Inbox className="h-4 w-4 mr-2" />
                Activity Feed
              </Button>
            </CardContent>
          </Card>

          {/* Pro Tip */}
          <Card className="bg-gradient-to-br from-primary/5 to-transparent border-primary/10">
            <CardHeader>
              <CardTitle className="text-base">Pro Tip</CardTitle>
            </CardHeader>
            <CardContent className="text-sm">
              <p>
                Pin projects to your sidebar for quick access from anywhere in
                MORIS
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default DashboardRoute;
