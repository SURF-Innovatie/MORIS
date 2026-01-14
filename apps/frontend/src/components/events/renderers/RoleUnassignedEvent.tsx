import { FC } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { ProjectEvent, ProjectEventType } from "@/api/events";
import { UserMinus, ShieldAlert } from "lucide-react";
import { EventMetaInfo } from "./EventMetaInfo";

export const RoleUnassignedEvent: FC<{ event: ProjectEvent }> = ({ event }) => {
  if (event.type !== ProjectEventType.ProjectRoleUnassigned || !event.person) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, email, avatarUrl, givenName, familyName } = event.person;
  const initials = (givenName?.[0] || "") + (familyName?.[0] || "");
  const roleName = event.projectRole?.name || "Member";

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center gap-3 opacity-75">
        <div className="relative">
          <Avatar className="h-10 w-10 border border-gray-200 grayscale">
            <AvatarImage src={avatarUrl} alt={name || "User"} />
            <AvatarFallback className="bg-gray-100 text-gray-500 font-medium">
              {initials || "U"}
            </AvatarFallback>
          </Avatar>
          <div className="absolute -bottom-1 -right-1 bg-white rounded-full p-0.5 shadow-sm border border-gray-100">
            <UserMinus className="h-3 w-3 text-red-500" />
          </div>
        </div>
        <div className="flex flex-col">
          <span className="text-sm font-medium text-gray-600 line-through decoration-red-400/50">
            {name}
          </span>
          <span className="text-xs text-gray-400">{email}</span>
        </div>
      </div>

      <div className="flex items-center gap-2 mt-1 px-3 py-1.5 bg-red-50/50 text-red-700 rounded-md border border-red-100/50 w-fit">
        <ShieldAlert className="h-4 w-4" />
        <span className="text-sm font-medium">Unassigned from: {roleName}</span>
      </div>

      <EventMetaInfo event={event} />
    </div>
  );
};
