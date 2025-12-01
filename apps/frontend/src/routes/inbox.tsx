import { useState } from "react";
import { format } from "date-fns";
import { Link } from "react-router-dom";
import { Bell, CheckCircle2, Inbox, CheckCheck, ClipboardCheck, Info } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useNotifications } from "@/context/NotificationContext";
import { ApprovalModal } from "@/components/project-edit/ApprovalModal";

const getIcon = (type?: string) => {
    switch (type) {
        case 'approval_request':
            return <ClipboardCheck className="h-5 w-5 text-blue-500" />;
        case 'status_update':
            return <Info className="h-5 w-5 text-green-500" />;
        default:
            return <Bell className="h-5 w-5 text-gray-500" />;
    }
};

const InboxRoute = () => {
    const { notifications, unreadCount, markAsRead } = useNotifications();
    const [selectedApproval, setSelectedApproval] = useState<{ projectId: string; eventId: string } | null>(null);

    const handleMarkAsRead = async (id: string) => {
        await markAsRead(id);
    };

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
                    {!notifications || notifications.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-12 text-center">
                            <Bell className="mb-4 h-12 w-12 text-muted-foreground/30" />
                            <p className="text-lg font-medium text-muted-foreground">
                                All caught up!
                            </p>
                            <p className="text-sm text-muted-foreground/70">
                                You have no new notifications.
                            </p>
                        </div>
                    ) : (
                        <div className="divide-y divide-border">
                            {notifications.map((notification) => {
                                const isApproval = notification.type === 'approval_request';

                                const Content = (
                                    <div
                                        className={`group flex items-start gap-4 p-4 transition-colors hover:bg-muted/50 ${!notification.read ? "bg-primary/5" : ""
                                            }`}
                                    >
                                        <div className="mt-1">
                                            {getIcon(notification.type)}
                                        </div>
                                        <div className="flex-1 space-y-1">
                                            <p className={`text-sm ${!notification.read ? "font-medium text-foreground" : "text-muted-foreground"}`}>
                                                {notification.message}
                                            </p>
                                            <p className="text-xs text-muted-foreground/70">
                                                {notification.sentAt
                                                    ? format(new Date(notification.sentAt), "MMM d, yyyy 'at' h:mm a")
                                                    : "Just now"}
                                            </p>
                                        </div>
                                        {!notification.read && (
                                            <Button
                                                variant="ghost"
                                                size="icon"
                                                className="h-8 w-8 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity"
                                                title="Mark as read"
                                                onClick={(e) => {
                                                    e.preventDefault();
                                                    e.stopPropagation();
                                                    handleMarkAsRead(notification.id!);
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
                                            onClick={() => setSelectedApproval({
                                                projectId: notification.projectId!,
                                                eventId: notification.eventId!
                                            })}
                                            className="cursor-pointer"
                                        >
                                            {Content}
                                        </div>
                                    );
                                }

                                const linkTarget = notification.projectId
                                    ? `/projects/${notification.projectId}`
                                    : null;

                                return (
                                    <div key={notification.id}>
                                        {linkTarget ? (
                                            <Link to={linkTarget} className="block">
                                                {Content}
                                            </Link>
                                        ) : (
                                            Content
                                        )}
                                    </div>
                                );
                            })}
                        </div>
                    )}
                </CardContent>
            </Card>

            {selectedApproval && (
                <ApprovalModal
                    isOpen={!!selectedApproval}
                    onClose={() => setSelectedApproval(null)}
                    projectId={selectedApproval.projectId}
                    eventId={selectedApproval.eventId}
                />
            )}
        </div>
    );
};

export default InboxRoute;
