import { useState, useEffect } from "react";
import {
  useGetOrganisationNodesSearch,
  useGetOrganisationNodesId,
} from "@api/moris";
import { OrganisationResponse } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Search, Check, Building2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { EMPTY_UUID } from "@/lib/constants";

// Helper to check if value is a valid, non-empty UUID
const isValidOrganisationId = (id: string | undefined): id is string => {
  return !!id && id !== EMPTY_UUID;
};

interface OrganisationSearchSelectProps {
  value?: string;
  onSelect: (
    organisationId: string,
    organisation: OrganisationResponse
  ) => void;
  disabled?: boolean;
  permission?: string;
}

export const OrganisationSearchSelect = ({
  value,
  onSelect,
  disabled,
  permission,
}: OrganisationSearchSelectProps) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [selectedOrganisation, setSelectedOrganisation] =
    useState<OrganisationResponse | null>(null);

  // Search organisations using the backend search endpoint
  const { data: results, isLoading } = useGetOrganisationNodesSearch(
    { q: query, permission },
    { query: { enabled: open && query.length >= 3 } }
  );

  // Fetch organisation by ID if value is provided but we don't have the selected organisation
  const { data: valueLookup } = useGetOrganisationNodesId(value || "", {
    query: {
      enabled: isValidOrganisationId(value) && !selectedOrganisation,
      staleTime: Infinity,
    },
  });

  // Sync selected organisation with value lookup
  useEffect(() => {
    if (isValidOrganisationId(value) && !selectedOrganisation && valueLookup) {
      setSelectedOrganisation(valueLookup);
    }
  }, [value, selectedOrganisation, valueLookup]);

  // Clear selection if value is externally set to undefined or empty
  useEffect(() => {
    if (!isValidOrganisationId(value) && selectedOrganisation) {
      setSelectedOrganisation(null);
    }
  }, [value]);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className="w-full justify-between"
          disabled={disabled}
        >
          {selectedOrganisation ? (
            <div className="flex items-center gap-2 overflow-hidden">
              <OrganisationAvatar organisation={selectedOrganisation} />
              <span className="truncate">{selectedOrganisation.name}</span>
            </div>
          ) : (
            <span className="text-muted-foreground">
              {isValidOrganisationId(value)
                ? "Loading..."
                : "Select organisation..."}
            </span>
          )}
          <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[350px] p-0"
        align="start"
        side="bottom"
        avoidCollisions={false}
      >
        <div className="flex items-center border-b px-3">
          <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input
            className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0 shadow-none"
            placeholder="Search organisations..."
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

          {!isLoading && query.length < 3 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              Type at least 3 characters to search...
            </div>
          )}

          {!isLoading && query.length >= 3 && results?.length === 0 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              No organisation found.
            </div>
          )}

          {!isLoading &&
            results?.map((organisation) => (
              <div
                key={organisation.id}
                className={cn(
                  "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                  value === organisation.id &&
                    "bg-accent text-accent-foreground"
                )}
                onClick={() => {
                  onSelect(organisation.id!, organisation);
                  setSelectedOrganisation(organisation);
                  setOpen(false);
                }}
              >
                <Check
                  className={cn(
                    "mr-2 h-4 w-4",
                    value === organisation.id ? "opacity-100" : "opacity-0"
                  )}
                />
                <div className="flex items-center gap-2 overflow-hidden">
                  <OrganisationAvatar organisation={organisation} />
                  <div className="flex flex-col">
                    <span className="font-medium">{organisation.name}</span>
                    {organisation.description && (
                      <span className="text-xs text-muted-foreground truncate max-w-[250px]">
                        {organisation.description}
                      </span>
                    )}
                  </div>
                </div>
              </div>
            ))}
        </div>
      </PopoverContent>
    </Popover>
  );
};

const OrganisationAvatar = ({
  organisation,
}: {
  organisation: OrganisationResponse;
}) => {
  if (organisation.avatarUrl) {
    return (
      <img
        src={organisation.avatarUrl}
        alt={organisation.name}
        className="h-6 w-6 rounded-full object-cover"
      />
    );
  }
  return (
    <Building2 className="h-6 w-6 p-1 rounded-full bg-secondary text-secondary-foreground" />
  );
};
