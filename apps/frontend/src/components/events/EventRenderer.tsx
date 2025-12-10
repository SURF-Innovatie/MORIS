import { FC } from "react";
import { Event } from "@/api/generated-orval/model";
import { PersonAddedEvent } from "./renderers/PersonAddedEvent";
import { PersonRemovedEvent } from "./renderers/PersonRemovedEvent";
import { ProductAddedEvent } from "./renderers/ProductAddedEvent";
import { ProductRemovedEvent } from "./renderers/ProductRemovedEvent";
import { DefaultEventRenderer } from "./renderers/DefaultEventRenderer";

interface EventRendererProps {
  event: Event;
  className?: string;
}

const RENDERER_REGISTRY: Record<string, FC<{ event: Event }>> = {
  "project.person_added": PersonAddedEvent,
  "project.person_removed": PersonRemovedEvent,
  "project.product_added": ProductAddedEvent,
  "project.product_removed": ProductRemovedEvent,
};

export const EventRenderer: FC<EventRendererProps> = ({ event, className }) => {
  const Renderer =
    (event.type && RENDERER_REGISTRY[event.type]) || DefaultEventRenderer;

  return (
    <div className={className}>
      <Renderer event={event} />
      {event.projectTitle && (
        <div className="mt-2 text-xs text-muted-foreground border-t pt-2 flex items-center gap-1">
          <span className="opacity-70">Project:</span>
          <span className="font-medium">{event.projectTitle}</span>
        </div>
      )}
    </div>
  );
};
