import { useState, useEffect } from "react";
import { useGetOrganisationNodesRorSearch } from "@api/moris";
import { RORItem } from "@api/model";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Search, Check } from "lucide-react";
import { cn } from "@/lib/utils";
import RorIcon from "@/components/icons/rorIcon";

interface RorSearchSelectProps {
  value?: string; // This is the ROR ID
  onSelect: (rorId: string, item: RORItem) => void;
  disabled?: boolean;
}

export const RorSearchSelect = ({
  value,
  onSelect,
  disabled,
}: RorSearchSelectProps) => {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const [selectedItem, setSelectedItem] = useState<RORItem | null>(null);

  // Using the hook with a manual query
  const { data: results, isLoading } = useGetOrganisationNodesRorSearch(
    { q: query },
    { query: { enabled: open && query.length > 2 } }
  );

  // Effect to fetch details if value is provided but we don't have the item (initial load)
  // We reuse the search endpoint which supports searching by ID
  const { data: valueLookup } = useGetOrganisationNodesRorSearch(
    { q: value || "" },
    {
      query: {
        enabled: !!value && !selectedItem,
        staleTime: Infinity, // Keep it cached
      },
    }
  );

  useEffect(() => {
    if (value && !selectedItem && valueLookup && valueLookup.length > 0) {
      // Try to find exact ID match first
      const match = valueLookup.find((i) => i.id === value) || valueLookup[0];
      if (match) {
        setSelectedItem(match);
      }
    }
  }, [value, selectedItem, valueLookup]);

  // If value changes to undefined externally, clear selection
  useEffect(() => {
    if (!value && selectedItem) {
      setSelectedItem(null);
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
          {selectedItem || (value && valueLookup?.[0]) ? (
            <div className="flex items-center gap-2 overflow-hidden">
              <RorIcon width={16} height={16} className="shrink-0" />
              <span className="truncate">
                {(selectedItem || valueLookup?.[0])?.name}{" "}
                <span className="text-muted-foreground font-normal">
                  (
                  {(selectedItem || valueLookup?.[0])?.id?.replace(
                    /^https?:\/\/ror\.org\//,
                    ""
                  )}
                  )
                </span>
              </span>
            </div>
          ) : (
            <span className="text-muted-foreground">
              {value ? value : "Select ROR organization..."}
            </span>
          )}
          <Search className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent
        className="w-[400px] p-0"
        align="start"
        side="bottom"
        avoidCollisions={false}
      >
        <div className="flex items-center border-b px-3">
          <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input
            className="flex h-11 w-full rounded-md bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground border-none focus-visible:ring-0 shadow-none"
            placeholder="Search ROR (min 3 chars)..."
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

          {!isLoading && query.length > 2 && results?.length === 0 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              No organization found.
            </div>
          )}

          {!isLoading && query.length <= 2 && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              Type at least 3 characters to search.
            </div>
          )}

          {!isLoading &&
            results?.map((item) => (
              <div
                key={item.id}
                className={cn(
                  "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none hover:bg-accent hover:text-accent-foreground",
                  value === item.id && "bg-accent text-accent-foreground"
                )}
                onClick={() => {
                  onSelect(item.id!, item);
                  setSelectedItem(item);
                  setOpen(false);
                }}
              >
                <Check
                  className={cn(
                    "mr-2 h-4 w-4",
                    value === item.id ? "opacity-100" : "opacity-0"
                  )}
                />
                <div className="flex items-center gap-2 overflow-hidden bg-transparent">
                  <div className="flex flex-col">
                    <span className="font-medium">
                      {item.name}{" "}
                      <span className="text-muted-foreground font-normal">
                        ({item.id?.replace(/^https?:\/\/ror\.org\//, "")})
                      </span>
                    </span>
                    {item.country && (
                      <span className="text-xs text-muted-foreground">
                        {item.addresses?.[0]?.city
                          ? `${item.addresses[0].city}, `
                          : ""}
                        {item.country.country_name}
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
