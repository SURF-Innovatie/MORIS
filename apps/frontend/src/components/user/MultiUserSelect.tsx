import { useState, useEffect } from "react";
import { useGetUsersSearch } from "@api/moris";
import { PersonResponse } from "@api/model";
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

interface MultiUserSelectProps {
  value: string[];
  onChange: (userIds: string[]) => void;
  disabled?: boolean;
  placeholder?: string;
  /** Optional initial person data for pre-selected users (used when editing) */
  initialPersons?: PersonResponse[];
}

export const MultiUserSelect = ({
  value,
  onChange,
  disabled,
  placeholder = "Select users...",
  initialPersons,
}: MultiUserSelectProps) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [selectedPersons, setSelectedPersons] = useState<PersonResponse[]>(
    initialPersons || []
  );

  // Sync initialPersons when they change (e.g., when editing a policy)
  useEffect(() => {
    if (initialPersons && initialPersons.length > 0) {
      setSelectedPersons(initialPersons);
    }
  }, [initialPersons]);

  const { data: results, isLoading } = useGetUsersSearch(
    { q: query },
    { query: { enabled: open && query.length > 0 } }
  );

  const handleSelect = (person: PersonResponse) => {
    if (!person.id) return;

    if (value.includes(person.id)) {
      // Remove if already selected
      onChange(value.filter((id) => id !== person.id));
      setSelectedPersons((prev) => prev.filter((p) => p.id !== person.id));
    } else {
      // Add to selection
      onChange([...value, person.id]);
      setSelectedPersons((prev) => [...prev, person]);
    }
  };

  const handleRemove = (personId: string) => {
    onChange(value.filter((id) => id !== personId));
    setSelectedPersons((prev) => prev.filter((p) => p.id !== personId));
  };

  // Show count for IDs without person data (fallback)
  const unknownCount = value.length - selectedPersons.length;

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
                {selectedPersons.map((person) => (
                  <Badge key={person.id} variant="secondary" className="mr-1">
                    {person.name || person.email}
                    <button
                      type="button"
                      className="ml-1 ring-offset-background rounded-full outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleRemove(person.id!);
                      }}
                    >
                      <X className="h-3 w-3" />
                    </button>
                  </Badge>
                ))}
                {unknownCount > 0 && (
                  <Badge variant="outline" className="mr-1">
                    +{unknownCount} user(s)
                  </Badge>
                )}
              </div>
            ) : (
              <span className="text-muted-foreground">{placeholder}</span>
            )}
            <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-[300px] p-0" align="start">
          <div className="flex items-center border-b px-3">
            <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
            <Input
              className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0"
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
              results?.map((person) => (
                <div
                  key={person.id}
                  className={cn(
                    "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                    value.includes(person.id!) &&
                      "bg-accent text-accent-foreground"
                  )}
                  onClick={() => handleSelect(person)}
                >
                  <div
                    className={cn(
                      "mr-2 h-4 w-4 flex items-center justify-center border rounded-sm",
                      value.includes(person.id!)
                        ? "bg-primary border-primary text-primary-foreground"
                        : "border-input"
                    )}
                  >
                    {value.includes(person.id!) && (
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
                    <UserAvatar person={person} />
                    <div className="flex flex-col">
                      <span>{person.name}</span>
                      <span className="text-xs text-muted-foreground">
                        {person.email}
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

const UserAvatar = ({ person }: { person: PersonResponse }) => {
  if (person.avatarUrl) {
    return (
      <img
        src={person.avatarUrl}
        alt={person.name}
        className="h-6 w-6 rounded-full object-cover"
      />
    );
  }
  return (
    <UserIcon className="h-6 w-6 p-1 rounded-full bg-secondary text-secondary-foreground" />
  );
};
