import { FC } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Event } from "@/api/generated-orval/model";

export const PersonAddedEvent: FC<{ event: Event }> = ({ event }) => {
  if (!event.person) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, email, avatar_url, orcid, givenName, familyName } =
    event.person;
  const initials = (givenName?.[0] || "") + (familyName?.[0] || "");

  return (
    <div className="flex items-center gap-3">
      <Avatar className="h-10 w-10 border border-gray-200">
        <AvatarImage src={avatar_url} alt={name || "User"} />
        <AvatarFallback className="bg-blue-100 text-blue-700 font-medium">
          {initials || "U"}
        </AvatarFallback>
      </Avatar>
      <div className="flex flex-col">
        <span className="text-sm font-medium text-gray-900">
          Person Added: {name}
        </span>
        <div className="flex items-center gap-2 text-xs text-gray-500">
          {email && <span>{email}</span>}
          {orcid && (
            <>
              <span className="text-gray-300">•</span>
              <span className="font-mono text-green-700">ORCID: {orcid}</span>
            </>
          )}
        </div>
      </div>
    </div>
  );
};
