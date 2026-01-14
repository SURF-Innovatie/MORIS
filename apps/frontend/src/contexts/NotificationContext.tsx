import React, { createContext, useContext } from "react";
import { useGetNotificationsMe, usePostNotificationsIdRead } from "@api/moris";
import { NotificationResponse } from "@api/model";

interface NotificationContextType {
  notifications: NotificationResponse[];
  unreadCount: number;
  isLoading: boolean;
  markAsRead: (id: string) => Promise<void>;
  refetch: () => void;
}

const NotificationContext = createContext<NotificationContextType | undefined>(
  undefined
);

export function NotificationProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const {
    data: notifications = [],
    isLoading,
    refetch,
  } = useGetNotificationsMe({
    query: {
      refetchInterval: 30000, // Poll every 30 seconds
    },
  });

  const { mutateAsync: markAsReadMutation } = usePostNotificationsIdRead();

  const unreadCount = notifications?.filter((n) => !n.read).length ?? 0;

  const markAsRead = async (id: string) => {
    try {
      await markAsReadMutation({ id });
      // Optimistically update or refetch
      refetch();
    } catch (error) {
      console.error("Failed to mark notification as read", error);
    }
  };

  React.useEffect(() => {
    const handleRefresh = () => {
      console.log("Refetching notifications due to event emission");
      refetch();
    };

    window.addEventListener("notifications:should-refresh", handleRefresh);

    return () => {
      window.removeEventListener("notifications:should-refresh", handleRefresh);
    };
  }, [refetch]);

  const value = {
    notifications,
    unreadCount,
    isLoading,
    markAsRead,
    refetch,
  };

  return (
    <NotificationContext.Provider value={value}>
      {children}
    </NotificationContext.Provider>
  );
}

export function useNotifications() {
  const context = useContext(NotificationContext);
  if (context === undefined) {
    throw new Error(
      "useNotifications must be used within a NotificationProvider"
    );
  }
  return context;
}
