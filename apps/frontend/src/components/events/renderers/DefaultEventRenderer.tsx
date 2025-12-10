import { FC } from "react";
import { Event } from "@/api/generated-orval/model";
import { CalendarDays, User } from "lucide-react";
import { format } from "date-fns";

export const DefaultEventRenderer: FC<{ event: Event }> = ({ event }) => {
  // Format "project.person_added" -> "Person Added"
  const formattedType =
    event.type
      ?.replace(/^project\./, "") // Remove project. prefix
      ?.replace(/_/g, " ") || "Event"; // Replace underscores with spaces

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center justify-between">
        <span className="text-sm font-medium text-gray-900 capitalize">
          {formattedType}
        </span>
      </div>
      <p className="text-sm text-gray-600">
        {event.details || "No details available"}
      </p>

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
            {/* In a real app we might resolve this UUID to a name if not available in event object */}
            <span className="font-mono text-[10px] opacity-70">
              {event.createdBy.substring(0, 8)}...
            </span>
          </div>
        )}
      </div>
    </div>
  );
};
