import { useState, useEffect, useMemo } from "react";
import { useGetUsersSearch, getGetPeopleIdQueryOptions } from "@api/moris";
import { UserPersonResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Badge } from "@/components/ui/badge";
import { Search, X, User as UserIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { useQueries } from "@tanstack/react-query";

interface MultiUserSelectProps {
  /** Array of person IDs (since UserPersonResponse.id is person_id) */
  value: string[];
  onChange: (personIds: string[]) => void;
  disabled?: boolean;
  placeholder?: string;
  /** Initial user data for pre-selected users (used when editing) */
  initialUsers?: UserPersonResponse[];
}

export const MultiUserSelect = ({
  value,
  onChange,
  disabled,
  placeholder = "Select users...",
  initialUsers,
}: MultiUserSelectProps) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [selectedUsers, setSelectedUsers] = useState<UserPersonResponse[]>(
    initialUsers || []
  );

  // Sync initialUsers when they change (e.g., when editing a policy)
  useEffect(() => {
    if (initialUsers && initialUsers.length > 0) {
      setSelectedUsers(initialUsers);
    }
  }, [initialUsers]);

  // Find IDs that don't have user data yet
  const missingIds = useMemo(() => {
    const selectedIds = new Set(selectedUsers.map((u) => u.id));
    return value.filter((id) => !selectedIds.has(id));
  }, [value, selectedUsers]);

  // Fetch missing person data using /people/:id endpoint
  const missingUserQueries = useQueries({
    queries: missingIds.map((id) => ({
      ...getGetPeopleIdQueryOptions(id),
      enabled: !!id,
      staleTime: Infinity,
    })),
  });

  // Sync fetched persons into selectedUsers
  useEffect(() => {
    const fetchedUsers = missingUserQueries
      .filter((q) => q.isSuccess && q.data)
      .map((q) => {
        const person = q.data!;
        // Map PersonResponse to UserPersonResponse format
        return {
          id: person.id,
          name: person.name,
          email: person.email,
          avatarUrl: person.avatarUrl,
          givenName: person.givenName,
          familyName: person.familyName,
          orcid: person.orcid,
        } as UserPersonResponse;
      });

    if (fetchedUsers.length > 0) {
      setSelectedUsers((prev) => {
        const existingIds = new Set(prev.map((u) => u.id));
        const newUsers = fetchedUsers.filter((u) => !existingIds.has(u.id));
        return newUsers.length > 0 ? [...prev, ...newUsers] : prev;
      });
    }
  }, [missingUserQueries.map((q) => q.data).join(",")]);

  const { data: results, isLoading } = useGetUsersSearch(
    { q: query },
    { query: { enabled: open && query.length > 0 } }
  );

  const handleSelect = (user: UserPersonResponse) => {
    if (!user.id) return;

    if (value.includes(user.id)) {
      // Remove if already selected
      onChange(value.filter((id) => id !== user.id));
      setSelectedUsers((prev) => prev.filter((u) => u.id !== user.id));
    } else {
      // Add to selection
      onChange([...value, user.id]);
      setSelectedUsers((prev) => [...prev, user]);
    }
  };

  const handleRemove = (userId: string) => {
    onChange(value.filter((id) => id !== userId));
    setSelectedUsers((prev) => prev.filter((u) => u.id !== userId));
  };

  // Check if we're still loading missing users
  const isLoadingMissingUsers = missingUserQueries.some((q) => q.isLoading);

  return (
    <div className="space-y-2">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            aria-expanded={open}
            className="w-full justify-between h-auto min-h-10"
            disabled={disabled}
          >
            {value.length > 0 ? (
              <div className="flex flex-wrap gap-1">
                {selectedUsers.map((user) => (
                  <Badge key={user.id} variant="secondary" className="mr-1">
                    {user.name || user.email}
                    <button
                      type="button"
                      className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleRemove(user.id!);
                      }}
                    >
                      <X className="h-3 w-3" />
                    </button>
                  </Badge>
                ))}
                {isLoadingMissingUsers && (
                  <Badge variant="outline" className="mr-1">
                    Loading...
                  </Badge>
                )}
              </div>
            ) : (
              <span className="text-muted-foreground">{placeholder}</span>
            )}
            <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent
          className="w-[300px] p-0"
          align="start"
          side="bottom"
          avoidCollisions={false}
        >
          <div className="flex items-center border-b px-3">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <Input
              className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0 shadow-none"
              placeholder="Search by name or email..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
            />
          </div>
          <div className="max-h-[300px] overflow-y-auto p-1">
            {isLoading && (
              <div className="py-6 text-center text-sm text-muted-foreground">
                Searching...
              </div>
            )}

            {!isLoading && query.length > 0 && results?.length === 0 && (
              <div className="py-6 text-center text-sm text-muted-foreground">
                No user found.
              </div>
            )}

            {!isLoading &&
              results?.map((user) => (
                <div
                  key={user.id}
                  className={cn(
                    "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                    value.includes(user.id!) &&
                      "bg-accent text-accent-foreground"
                  )}
                  onClick={() => handleSelect(user)}
                >
                  <div
                    className={cn(
                      "mr-2 h-4 w-4 flex items-center justify-center border rounded-sm",
                      value.includes(user.id!)
                        ? "bg-primary border-primary text-primary-foreground"
                        : "border-input"
                    )}
                  >
                    {value.includes(user.id!) && (
                      <svg
                        className="h-3 w-3"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                    )}
                  </div>
                  <div className="flex items-center gap-2">
                    <UserAvatar user={user} />
                    <div className="flex flex-col">
                      <span>{user.name}</span>
                      <span className="text-xs text-muted-foreground">
                        {user.email}
                      </span>
                    </div>
                  </div>
                </div>
              ))}
          </div>
        </PopoverContent>
      </Popover>

      {/* Show selected users count */}
      {value.length > 0 && (
        <div className="flex flex-wrap gap-1 text-xs text-muted-foreground">
          {value.length} user{value.length > 1 ? "s" : ""} selected
        </div>
      )}
    </div>
  );
};

const UserAvatar = ({ user }: { user: UserPersonResponse }) => {
  if (user.avatarUrl) {
    return (
      <img
        src={user.avatarUrl}
        alt={user.name}
        className="h-6 w-6 rounded-full object-cover"
      />
    );
  }
  return (
    <UserIcon className="h-6 w-6 p-1 rounded-full bg-secondary text-secondary-foreground" />
  );
};
