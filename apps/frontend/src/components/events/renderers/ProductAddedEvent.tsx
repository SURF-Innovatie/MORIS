import { FC } from "react";
import { Badge } from "@/components/ui/badge";
import { ProjectEventType } from "@/api/events";
import { FileText, Plus } from "lucide-react";
import { EventMetaInfo } from "./EventMetaInfo";
import { EventRendererBaseProps } from "../types";

export const ProductAddedEvent: FC<EventRendererBaseProps> = ({
  event,
  variant = "normal",
}) => {
  if (event.type !== ProjectEventType.ProductAdded || !event.product) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, type, doi } = event.product;

  if (variant === "compact") {
    return (
      <div className="flex items-center gap-2">
        <div className="p-1 bg-blue-50 text-blue-600 rounded border border-blue-100">
          <Plus className="h-3.5 w-3.5" />
        </div>
        <div className="flex items-center gap-1.5 min-w-0 flex-1">
          <span
            className="text-sm font-medium text-gray-900 truncate max-w-[200px]"
            title={name}
          >
            {name}
          </span>
          {type && (
            <Badge
              variant="secondary"
              className="text-[10px] px-1 h-4 shrink-0"
            >
              {type}
            </Badge>
          )}
          <span className="text-gray-300">•</span>
          <EventMetaInfo event={event} variant="compact" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-start gap-3">
        <div className="mt-1 p-2 bg-blue-50 text-blue-600 rounded-lg border border-blue-100">
          <FileText className="h-5 w-5" />
        </div>
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="text-sm font-medium text-gray-900">
              Product Added: {name}
            </span>
            {type && (
              <Badge variant="secondary" className="text-[10px] px-1.5 h-5">
                {type}
              </Badge>
            )}
          </div>
          {doi && (
            <a
              href={`https://doi.org/${doi}`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-xs text-blue-600 hover:underline font-mono truncate max-w-[200px] sm:max-w-none"
            >
              DOI: {doi}
            </a>
          )}
        </div>
      </div>
      <EventMetaInfo event={event} />
    </div>
  );
};
