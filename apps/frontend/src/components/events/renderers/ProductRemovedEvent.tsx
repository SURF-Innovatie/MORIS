import { FC } from "react";
import { ProjectEvent, ProjectEventType } from "@/api/events";
import { FileMinus } from "lucide-react";
import { EventMetaInfo } from "./EventMetaInfo";

export const ProductRemovedEvent: FC<{ event: ProjectEvent }> = ({ event }) => {
  if (event.type !== ProjectEventType.ProductRemoved || !event.product) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, type } = event.product;

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center gap-3 opacity-75">
        <div className="p-2 bg-gray-50 text-gray-400 rounded-lg border border-gray-100 relative">
          <FileMinus className="h-5 w-5" />
        </div>
        <div className="flex flex-col">
          <span className="text-sm font-medium text-gray-600 line-through decoration-red-400/50">
            Product Removed: {name}
          </span>
          {type && (
            <span className="text-xs text-gray-400 lowercase italic">
              {type}
            </span>
          )}
        </div>
      </div>
      <EventMetaInfo event={event} />
    </div>
  );
};
