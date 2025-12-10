import { FC } from "react";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Event } from "@/api/generated-orval/model";
import { UserMinus } from "lucide-react";

export const PersonRemovedEvent: FC<{ event: Event }> = ({ event }) => {
  if (!event.person) {
    return <div className="text-sm text-gray-600">{event.details}</div>;
  }

  const { name, email, avatar_url, givenName, familyName } = event.person;
  const initials = (givenName?.[0] || "") + (familyName?.[0] || "");

  return (
    <div className="flex items-center gap-3 opacity-75">
      <div className="relative">
        <Avatar className="h-10 w-10 border border-gray-200 grayscale">
          <AvatarImage src={avatar_url} alt={name || "User"} />
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
          Person Removed: {name}
        </span>
        <span className="text-xs text-gray-400">{email}</span>
      </div>
    </div>
  );
};
