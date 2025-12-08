import { CheckCheck, Inbox } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useNotifications } from "@/contexts/NotificationContext";
import { NotificationList } from "@/components/notifications/NotificationList";

const InboxRoute = () => {
    const { notifications, unreadCount, markAsRead } = useNotifications();

    const handleMarkAllAsRead = async () => {
        if (!notifications) return;
        const unreadNotifications = notifications.filter((n) => !n.read);
        await Promise.all(unreadNotifications.map((n) => markAsRead(n.id!)));
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                    <Inbox className="h-6 w-6 text-muted-foreground" />
                    <h1 className="text-2xl font-semibold tracking-tight">Inbox</h1>
                    {unreadCount > 0 && (
                        <Badge variant="secondary" className="ml-2">
                            {unreadCount} Unread
                        </Badge>
                    )}
                </div>
                {unreadCount > 0 && (
                    <Button variant="outline" size="sm" onClick={handleMarkAllAsRead}>
                        <CheckCheck className="mr-2 h-4 w-4" />
                        Mark all as read
                    </Button>
                )}
            </div>

            <Card>
                <CardContent className="p-0">
                    <NotificationList />
                </CardContent>
            </Card>
        </div>
    );
};

export default InboxRoute;
