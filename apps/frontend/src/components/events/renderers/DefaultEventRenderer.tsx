import { FC } from "react";
import { EventMetaInfo } from "./EventMetaInfo";
import { EventRendererBaseProps } from "../types";
import { Activity } from "lucide-react";

export const DefaultEventRenderer: FC<EventRendererBaseProps> = ({
  event,
  variant = "normal",
}) => {
  // Format "project.person_added" -> "Person Added"
  const formattedType =
    event.friendlyName ||
    event.type
      ?.replace(/^project\./, "") // Remove project. prefix
      ?.replace(/_/g, " ") ||
    "Event"; // Replace underscores with spaces

  if (variant === "compact") {
    return (
      <div className="flex items-center gap-2">
        <div className="p-1 bg-gray-50 text-gray-500 rounded border border-gray-100">
          <Activity className="h-3.5 w-3.5" />
        </div>
        <div className="flex items-center gap-1.5 min-w-0 flex-1">
          <span className="text-sm font-medium text-gray-900 capitalize truncate">
            {formattedType}
          </span>
          <span className="text-gray-300">•</span>
          <EventMetaInfo event={event} variant="compact" />
        </div>
      </div>
    );
  }

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

      <EventMetaInfo event={event} />
    </div>
  );
};
