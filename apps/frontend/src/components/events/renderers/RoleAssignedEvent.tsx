import { FC } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ProjectEventType } from "@/api/events";
import { EventMetaInfo } from "./EventMetaInfo";
import { ShieldCheck } from "lucide-react";
import { EventRendererBaseProps } from "../types";

export const RoleAssignedEvent: FC<EventRendererBaseProps> = ({
  event,
  variant = "normal",
}) => {
  if (event.type !== ProjectEventType.ProjectRoleAssigned || !event.person) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, email, avatarUrl, givenName, familyName } = event.person;
  const initials = (givenName?.[0] || "") + (familyName?.[0] || "");
  const roleName = event.projectRole?.name || "Member";

  if (variant === "compact") {
    return (
      <div className="flex items-center gap-2">
        <Avatar className="h-6 w-6 border border-gray-200">
          <AvatarImage src={avatarUrl} alt={name || "User"} />
          <AvatarFallback className="bg-blue-100 text-blue-700 text-xs font-medium">
            {initials || "U"}
          </AvatarFallback>
        </Avatar>
        <div className="flex items-center gap-1.5 min-w-0 flex-1">
          <span className="text-sm font-medium text-gray-900 truncate">
            {name}
          </span>
          <ShieldCheck className="h-3.5 w-3.5 text-blue-600 shrink-0" />
          <span className="text-sm text-blue-600 font-medium">{roleName}</span>
          <span className="text-gray-300">•</span>
          <EventMetaInfo event={event} variant="compact" />
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center gap-3">
        <Avatar className="h-10 w-10 border border-gray-200">
          <AvatarImage src={avatarUrl} alt={name || "User"} />
          <AvatarFallback className="bg-blue-100 text-blue-700 font-medium">
            {initials || "U"}
          </AvatarFallback>
        </Avatar>
        <div className="flex flex-col">
          <span className="text-sm font-medium text-gray-900">{name}</span>
          <span className="text-xs text-gray-500">{email}</span>
        </div>
      </div>

      <div className="flex items-center gap-2 mt-1 px-3 py-1.5 bg-blue-50/50 text-blue-700 rounded-md border border-blue-100/50 w-fit">
        <ShieldCheck className="h-4 w-4" />
        <span className="text-sm font-medium">Assigned as: {roleName}</span>
      </div>

      <EventMetaInfo event={event} />
    </div>
  );
};
