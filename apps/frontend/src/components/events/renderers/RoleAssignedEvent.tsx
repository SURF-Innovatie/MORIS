import { FC } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ProjectEvent, ProjectEventType } from "@/api/events";
import { EventMetaInfo } from "./EventMetaInfo";
import { ShieldCheck } from "lucide-react";

export const RoleAssignedEvent: FC<{ event: ProjectEvent }> = ({ event }) => {
  if (event.type !== ProjectEventType.ProjectRoleAssigned || !event.person) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, email, avatarUrl, givenName, familyName } = event.person;
  const initials = (givenName?.[0] || "") + (familyName?.[0] || "");
  const roleName = event.projectRole?.name || "Member";

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
