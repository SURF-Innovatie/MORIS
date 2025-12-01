import { useState } from "react";
import { Bell } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { Button } from "./ui/button";
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "./ui/popover";
import { Badge } from "./ui/badge";
import { useGetNotifications, usePutNotificationsIdRead } from "@api/moris";
import { NotificationResponse } from "@/api/generated-orval/model";

export const NotificationBell = () => {
    const [open, setOpen] = useState(false);
    const navigate = useNavigate();
    const { data: notifications, refetch } = useGetNotifications({
        query: {
            refetchInterval: 30000, // Poll every 30 seconds
        }
    });
    const { mutate: markAsRead } = usePutNotificationsIdRead();

    const unreadCount = notifications?.filter((n) => !n.read).length || 0;

    const handleNotificationClick = (notification: NotificationResponse) => {
        if (!notification.read) {
            markAsRead({ id: notification.id! }, {
                onSuccess: () => {
                    refetch();
                }
            });
        }

        // Navigate if project ID is present
        const projectId = notification.projectId === "00000000-0000-0000-0000-000000000000" ? undefined : notification.projectId;
        if (projectId) {
            navigate(`/projects/${projectId}`);
            setOpen(false);
        }
    };

    const handleMarkAsRead = (e: React.MouseEvent, id: string) => {
        e.stopPropagation(); // Prevent navigation when clicking just "Mark as read"
        markAsRead({ id }, {
            onSuccess: () => {
                refetch();
            }
        });
    };

    return (
        <Popover open={open} onOpenChange={setOpen}>
            <PopoverTrigger asChild>
                <Button variant="ghost" size="icon" className="relative">
                    <Bell className="h-5 w-5" />
                    {unreadCount > 0 && (
                        <Badge
                            variant="destructive"
                            className="absolute -top-1 -right-1 h-5 w-5 flex items-center justify-center p-0 text-xs rounded-full"
                        >
                            {unreadCount}
                        </Badge>
                    )}
                </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80 p-0" align="end">
                <div className="p-4 border-b border-border">
                    <h4 className="font-semibold leading-none">Notifications</h4>
                </div>
                <div className="max-h-[300px] overflow-y-auto">
                    {notifications?.length === 0 ? (
                        <div className="p-4 text-center text-sm text-muted-foreground">
                            No notifications
                        </div>
                    ) : (
                        <div className="flex flex-col">
                            {notifications?.map((notification) => (
                                <div
                                    key={notification.id}
                                    onClick={() => handleNotificationClick(notification)}
                                    className={`p-4 border-b border-border last:border-0 hover:bg-muted/50 transition-colors cursor-pointer ${!notification.read ? "bg-muted/20" : ""
                                        }`}
                                >
                                    <div className="flex justify-between items-start gap-2">
                                        <p className="text-sm font-medium leading-none">
                                            {notification.message}
                                        </p>
                                        {!notification.read && (
                                            <Button
                                                variant="ghost"
                                                size="sm"
                                                className="h-auto p-0 text-xs text-primary hover:text-primary/80"
                                                onClick={(e) => handleMarkAsRead(e, notification.id!)}
                                            >
                                                Mark as read
                                            </Button>
                                        )}
                                    </div>
                                    <span className="text-xs text-muted-foreground mt-1 block">
                                        {new Date(notification.sentAt!).toLocaleDateString()}
                                    </span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </PopoverContent>
        </Popover>
    );
};
