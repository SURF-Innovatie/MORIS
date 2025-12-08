import { useNavigate } from "react-router-dom";
import { Inbox } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
} from "@/components/ui/card";
import { useNotifications } from "@/contexts/NotificationContext";
import { ProjectList } from "@/components/projects/ProjectList";
import { NotificationList } from "@/components/notifications/NotificationList";

const DashboardRoute = () => {
  const navigate = useNavigate();
  const { unreadCount } = useNotifications();

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

          <Card className="border-black/10shadow-none bg-transparent">
            <CardContent className="p-0">
              <NotificationList limit={5} />
            </CardContent>
          </Card>
        </section>
      </div>
    </div>
  );
};

export default DashboardRoute;
