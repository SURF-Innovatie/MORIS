import { FC } from "react";

import { RoleAssignedEvent } from "./renderers/RoleAssignedEvent";
import { RoleUnassignedEvent } from "./renderers/RoleUnassignedEvent";
import { ProductAddedEvent } from "./renderers/ProductAddedEvent";
import { ProductRemovedEvent } from "./renderers/ProductRemovedEvent";
import { DefaultEventRenderer } from "./renderers/DefaultEventRenderer";

import { ProjectEvent, ProjectEventType } from "@/api/events";

interface EventRendererProps {
  event: ProjectEvent; // Using the strict discriminated union
  className?: string;
}

const RENDERER_REGISTRY: Partial<
  Record<ProjectEventType, FC<{ event: ProjectEvent }>>
> = {
  [ProjectEventType.ProjectRoleAssigned]: RoleAssignedEvent,
  [ProjectEventType.ProjectRoleUnassigned]: RoleUnassignedEvent,
  [ProjectEventType.ProductAdded]: ProductAddedEvent,
  [ProjectEventType.ProductRemoved]: ProductRemovedEvent,
};

export const EventRenderer: FC<EventRendererProps> = ({ event, className }) => {
  const Renderer =
    RENDERER_REGISTRY[event.type as ProjectEventType] || DefaultEventRenderer;

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
