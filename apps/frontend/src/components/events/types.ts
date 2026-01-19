import { ProjectEvent } from "@/api/events";

export type EventDisplayVariant = "normal" | "compact";

export interface EventRendererBaseProps {
  event: ProjectEvent;
  variant?: EventDisplayVariant;
}
