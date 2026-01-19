import { FC } from "react";
import { format } from "date-fns";
import { CalendarDays, User } from "lucide-react";
import { EventDisplayVariant } from "../types";
import { ProjectEvent } from "@/api/events";

interface EventMetaInfoProps {
  event: ProjectEvent;
  variant?: EventDisplayVariant;
}

export const EventMetaInfo: FC<EventMetaInfoProps> = ({
  event,
  variant = "normal",
}) => {
  if (variant === "compact") {
    return (
      <div className="flex items-center gap-2 text-xs text-gray-400">
        {event.at && <span>{format(new Date(event.at), "MMM d, h:mm a")}</span>}
        {event.creator?.name && (
          <>
            <span>•</span>
            <span>{event.creator.name}</span>
          </>
        )}
      </div>
    );
  }

  return (
    <div className="flex items-center gap-4 mt-2 text-xs text-gray-400">
      {event.at && (
        <div className="flex items-center gap-1">
          <CalendarDays className="h-3 w-3" />
          <span>{format(new Date(event.at), "MMM d, yyyy h:mm a")}</span>
        </div>
      )}
      {event.createdBy && (
        <div className="flex items-center gap-1">
          <User className="h-3 w-3" />
          <span className="font-mono text-[10px] opacity-70">
            {event.creator?.name || `${event.createdBy.substring(0, 8)}...`}
          </span>
        </div>
      )}
    </div>
  );
};
