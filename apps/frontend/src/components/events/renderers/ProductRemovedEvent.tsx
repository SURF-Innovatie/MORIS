import { FC } from "react";
import { ProjectEventType } from "@/api/events";
import { FileMinus, Minus } from "lucide-react";
import { EventMetaInfo } from "./EventMetaInfo";
import { EventRendererBaseProps } from "../types";

export const ProductRemovedEvent: FC<EventRendererBaseProps> = ({
  event,
  variant = "normal",
}) => {
  if (event.type !== ProjectEventType.ProductRemoved || !event.product) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, type } = event.product;

  if (variant === "compact") {
    return (
      <div className="flex items-center gap-2 opacity-75">
        <div className="p-1 bg-gray-50 text-gray-400 rounded border border-gray-100">
          <Minus className="h-3.5 w-3.5" />
        </div>
        <div className="flex items-center gap-1.5 min-w-0 flex-1">
          <span className="text-sm text-gray-600 line-through decoration-red-400/50 truncate">
            {name}
          </span>
          {type && (
            <span className="text-xs text-gray-400 italic shrink-0">
              {type}
            </span>
          )}
          <span className="text-gray-300">•</span>
          <EventMetaInfo event={event} variant="compact" />
        </div>
      </div>
    );
  }

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
