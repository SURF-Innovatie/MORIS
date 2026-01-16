import { FC } from "react";
import { ProjectEvent } from "@/api/events";
import { EventMetaInfo } from "./EventMetaInfo";

export const DefaultEventRenderer: FC<{ event: ProjectEvent }> = ({
  event,
}) => {
  // Format "project.person_added" -> "Person Added"
  const formattedType =
    event.friendlyName ||
    event.type
      ?.replace(/^project\./, "") // Remove project. prefix
      ?.replace(/_/g, " ") ||
    "Event"; // Replace underscores with spaces

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
